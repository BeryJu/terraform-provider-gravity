package provider

import (
	"bytes"
	"io"
	"log"
	"net/http"

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
