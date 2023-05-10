package provider

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setWrapper(d *schema.ResourceData, key string, data interface{}) {
	err := d.Set(key, data)
	if err != nil {
		panic(err)
	}
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
	_, er := io.Copy(buff, r.Body)
	if er != nil {
		log.Printf("[DEBUG] Gravity: failed to read response: %s", er.Error())
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

func tfListMap(raw []interface{}) []map[string]string {
	values := make([]map[string]string, len(raw))
	for i, rh := range raw {
		values[i] = tfMap(rh.(map[string]interface{}))
	}
	return values
}

// StringHashcode hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func StringHashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
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
