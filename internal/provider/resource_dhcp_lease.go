package provider

import (
	"context"
	"fmt"
	"strings"

	"beryju.io/gravity/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDHCPLease() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDHCPLeaseCreate,
		ReadContext:   resourceDHCPLeaseRead,
		UpdateContext: resourceDHCPLeaseUpdate,
		DeleteContext: resourceDHCPLeaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDHCPLeaseImport,
		},
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reservation": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_lease_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Lease time offered to this client, as a Go duration (e.g. `12h`). Overrides the scope's `lease_ttl` when set.",
				ValidateDiagFunc: validateDuration,
			},
			"dns_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS zone this lease is published into. Overrides the scope's `dns.zone` when set.",
			},
			"vendor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vendor reported by the client.",
			},
		},
	}
}

func resourceDHCPLeaseSchemaToModel(d *schema.ResourceData) *api.DhcpAPILeasesPutInput {
	m := api.DhcpAPILeasesPutInput{
		Address:  d.Get("address").(string),
		Hostname: d.Get("hostname").(string),
	}
	if res := d.Get("reservation").(bool); res {
		m.Expiry = api.PtrInt64(-1)
	}
	if v, ok := d.GetOk("description"); ok {
		m.Description = api.PtrString(v.(string))
	}
	if v, ok := d.GetOk("address_lease_time"); ok {
		m.AddressLeaseTime = v.(string)
	}
	if v, ok := d.GetOk("dns_zone"); ok {
		m.DnsZone = api.PtrString(v.(string))
	}
	return &m
}

func resourceDHCPLeaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	req := resourceDHCPLeaseSchemaToModel(d)
	scope := d.Get("scope").(string)
	identifier := d.Get("identifier").(string)

	hr, err := c.client.RolesDhcpAPI.DhcpPutLeases(ctx).
		Scope(scope).
		Identifier(identifier).
		DhcpAPILeasesPutInput(*req).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	d.SetId(identifier)
	return resourceDHCPLeaseRead(ctx, d, m)
}

func resourceDHCPLeaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	scope := d.Get("scope").(string)
	identifier := d.Get("identifier").(string)

	res, hr, err := c.client.RolesDhcpAPI.DhcpGetLeases(ctx).
		Scope(scope).
		Identifier(identifier).
		Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Leases) < 1 {
		d.SetId("")
		return diag.Diagnostics{}
	}
	l := res.Leases[0]
	d.SetId(l.Identifier)
	setWrapper(d, "identifier", l.Identifier)
	setWrapper(d, "address", l.Address)
	setWrapper(d, "hostname", l.Hostname)
	setWrapper(d, "scope", l.ScopeKey)
	setWrapper(d, "description", l.Description)
	setWrapper(d, "address_lease_time", l.AddressLeaseTime)
	setWrapper(d, "dns_zone", stringValue(l.DnsZone))
	// A reservation is encoded as a negative expiry. Leases that expire
	// normally carry a timestamp, and the field is omitted entirely for
	// leases the server has not assigned an expiry to.
	setWrapper(d, "reservation", l.Expiry != nil && *l.Expiry <= -1)
	if l.Info != nil {
		setWrapper(d, "vendor", stringValue(l.Info.Vendor))
	} else {
		setWrapper(d, "vendor", "")
	}
	return diags
}

func resourceDHCPLeaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := resourceDHCPLeaseCreate(ctx, d, m)
	if diag != nil {
		return diag
	}
	return resourceDHCPLeaseRead(ctx, d, m)
}

func resourceDHCPLeaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	scope := d.Get("scope").(string)
	identifier := d.Get("identifier").(string)

	hr, err := c.client.RolesDhcpAPI.DhcpDeleteLeases(ctx).
		Scope(scope).
		Identifier(identifier).
		Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}

// resourceDHCPLeaseImport accepts "<scope>/<identifier>", since the scope is
// needed to look the lease up but is not recoverable from the identifier.
func resourceDHCPLeaseImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	scope, identifier, ok := strings.Cut(d.Id(), "/")
	if !ok {
		return nil, fmt.Errorf("invalid ID %q, expected \"<scope>/<identifier>\"", d.Id())
	}
	setWrapper(d, "scope", scope)
	setWrapper(d, "identifier", identifier)
	d.SetId(identifier)
	return []*schema.ResourceData{d}, nil
}
