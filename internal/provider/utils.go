package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func setWrapper(d *schema.ResourceData, key string, data interface{}) {
	err := d.Set(key, data)
	if err != nil {
		panic(err)
	}
}

func StringInEnum[T ~string](items []T) schema.SchemaValidateDiagFunc {
	nv := make([]string, len(items))
	for i, v := range items {
		nv[i] = string(v)
	}
	return validation.ToDiagFunc(validation.StringInSlice(nv, false))
}

func EnumToDescription[T ~string](allowed []T) string {
	sb := &strings.Builder{}
	sb.WriteString("Allowed values:\n")
	for _, v := range allowed {
		fmt.Fprintf(sb, "  - `%s`\n", v)
	}
	return sb.String()
}

func sliceToString(in []interface{}) []string {
	sl := make([]string, len(in))
	for i, m := range in {
		sl[i] = m.(string)
	}
	return sl
}

func httpToDiag(d *schema.ResourceData, r *http.Response, err error) diag.Diagnostics {
	if r == nil {
		return diag.Errorf("HTTP Error '%s' without http response", err.Error())
	}
	if r.StatusCode == 404 {
		d.SetId("")
		return diag.Diagnostics{}
	}
	buff := &bytes.Buffer{}
	if r.Body != nil {
		_, er := io.Copy(buff, r.Body)
		if er != nil {
			log.Printf("[DEBUG] Gravity: failed to read response: %s", er.Error())
		}
	}
	log.Printf("[DEBUG] Gravity: error response: %s", buff.String())
	return diag.Errorf("HTTP Error '%s' during request '%s %s': \"%s\"", err.Error(), r.Request.Method, r.Request.URL.Path, buff.String())
}

func tfMap(raw map[string]interface{}) map[string]string {
	x := make(map[string]string)
	for k, v := range raw {
		x[k] = v.(string)
	}
	return x
}

// The API marshals optional fields as pointers and omits them entirely when
// unset. Terraform has no notion of an unset primitive, so an absent value is
// read back as the type's zero value.

func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func boolValue(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func int32Value(v *int32) int {
	if v == nil {
		return 0
	}
	return int(*v)
}

func int64Value(v *int64) int {
	if v == nil {
		return 0
	}
	return int(*v)
}

// stringSlice normalises a possibly-nil API slice into an empty slice, so that
// state matches a config that simply omitted the attribute.
func stringSlice(v []string) []string {
	if v == nil {
		return []string{}
	}
	return v
}

// marshalJSONSlice encodes a slice as JSON, rendering a nil slice as `[]`
// rather than `null` so it round-trips against a `jsonencode([])` config.
func marshalJSONSlice[T any](in []T) (string, error) {
	if in == nil {
		in = []T{}
	}
	b, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// validateDuration checks a value parses as a Go duration, e.g. "12h30m".
func validateDuration(v any, p cty.Path) diag.Diagnostics {
	if _, err := time.ParseDuration(v.(string)); err != nil {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Invalid duration",
			Detail:   fmt.Sprintf("%q is not a valid duration: %s", v.(string), err),
		}}
	}
	return nil
}

func validateMustBeLowercase(summary string) schema.SchemaValidateDiagFunc {
	return func(v any, p cty.Path) diag.Diagnostics {
		value := v.(string)
		var diags diag.Diagnostics
		if strings.ToLower(value) != value {
			diag := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  summary,
				Detail:   fmt.Sprintf("%q is not lowercase", value),
			}
			diags = append(diags, diag)
		}
		return diags
	}
}
