package provider

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Writing a role's configuration makes the instance restart that role, and the
// new configuration is only swapped in once it comes back up. Two things
// follow from that, both of which these resources have to absorb:
//
//   - Restarting the API role tears down the listener this provider talks to,
//     so a request can fail at the transport level, or be sent over a pooled
//     connection belonging to a listener that no longer exists.
//   - A read taken immediately after a write can still return the previous
//     configuration, so writes poll until the server reports what they asked
//     for.
const (
	roleConfigAttempts = 20
	roleConfigBackoff  = 250 * time.Millisecond
)

// roleConfigResource builds a resource for one of Gravity's `/roles/<name>`
// endpoints. A role's configuration is a singleton that always exists, so
// these resources adopt the existing configuration rather than creating one,
// and destroying one only drops it from state.
//
// The endpoint replaces the stored configuration wholesale rather than merging
// into it, so every attribute is sent on every write. Attributes therefore
// carry defaults matching the server's own, so that omitting one does not
// quietly reset the role to a zero value.
type roleConfigResource[C any] struct {
	name        string
	description string
	schema      map[string]*schema.Schema
	get         func(ctx context.Context, c *APIClient) (C, *http.Response, error)
	put         func(ctx context.Context, c *APIClient, cfg C) (*http.Response, error)
	toModel     func(d *schema.ResourceData) C
	toState     func(d *schema.ResourceData, cfg C)
	// settled reports whether a configuration read back from the server
	// reflects one that was sent to it. Defaults to a deep comparison, which
	// suits roles whose every field round-trips unchanged.
	settled func(sent C, got C) bool
}

func (r roleConfigResource[C]) resource() *schema.Resource {
	return &schema.Resource{
		Description:   r.description,
		CreateContext: r.write,
		ReadContext:   r.read,
		UpdateContext: r.write,
		DeleteContext: r.delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: r.schema,
	}
}

func (r roleConfigResource[C]) write(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	sent := r.toModel(d)

	hr, err := r.putWithRetry(ctx, c, sent)
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	d.SetId(r.name)

	settled := r.settled
	if settled == nil {
		settled = func(sent C, got C) bool { return reflect.DeepEqual(sent, got) }
	}

	var got C
	for attempt := 0; ; attempt++ {
		got, hr, err = r.get(ctx, c)
		if err == nil && settled(sent, got) {
			break
		}
		// Anything the server actually answered is a real error, whereas a
		// failure with no response at all is the listener restarting.
		if err != nil && hr != nil {
			return httpToDiag(d, hr, err)
		}
		if attempt == roleConfigAttempts-1 {
			if err != nil {
				return httpToDiag(d, hr, err)
			}
			// The write was accepted but the server still reports something
			// else. Record what it has; the next plan shows the difference.
			break
		}
		if err := r.backoff(ctx); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(r.name)
	r.toState(d, got)
	return diag.Diagnostics{}
}

// putWithRetry resends a request that never reached the server. Replacing a
// role's configuration is idempotent, so this is safe, and it is what happens
// when the request goes out over a connection left pooled from a listener that
// has since restarted.
func (r roleConfigResource[C]) putWithRetry(ctx context.Context, c *APIClient, cfg C) (*http.Response, error) {
	var hr *http.Response
	var err error
	for attempt := 0; ; attempt++ {
		hr, err = r.put(ctx, c, cfg)
		if err == nil || !r.worthRetrying(hr, err, attempt) {
			return hr, err
		}
		if err := r.backoff(ctx); err != nil {
			return hr, err
		}
	}
}

func (r roleConfigResource[C]) read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	for attempt := 0; ; attempt++ {
		cfg, hr, err := r.get(ctx, c)
		if err == nil {
			d.SetId(r.name)
			r.toState(d, cfg)
			return diag.Diagnostics{}
		}
		if !r.worthRetrying(hr, err, attempt) {
			return httpToDiag(d, hr, err)
		}
		if err := r.backoff(ctx); err != nil {
			return diag.FromErr(err)
		}
	}
}

// worthRetrying reports whether a failure looks like the listener being
// restarted rather than a decision the server made. Anything the server
// actually answered is a real error and is reported as-is.
func (r roleConfigResource[C]) worthRetrying(hr *http.Response, err error, attempt int) bool {
	return hr == nil && attempt < roleConfigAttempts-1
}

func (r roleConfigResource[C]) backoff(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(roleConfigBackoff):
		return nil
	}
}

// delete forgets the role configuration. The role itself cannot be removed and
// its configuration keeps whatever values were last applied.
func (r roleConfigResource[C]) delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{}
}
