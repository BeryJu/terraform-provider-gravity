package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"beryju.io/gravity/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTFTPFile() *schema.Resource {
	return &schema.Resource{
		Description:   "File served by the TFTP role.",
		CreateContext: resourceTFTPFileCreate,
		ReadContext:   resourceTFTPFileRead,
		UpdateContext: resourceTFTPFileUpdate,
		DeleteContext: resourceTFTPFileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTFTPFileImport,
		},
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Host the file is served for. Use `*` to serve it to every " +
					"client regardless of the address they connect to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Path the file is served under.",
			},
			"content": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"content", "content_base64"},
				Description:  "File contents as text. Use `content_base64` for binary files.",
			},
			"content_base64": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"content", "content_base64"},
				Description:  "File contents, base64-encoded.",
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// resourceTFTPFileData returns the file contents in the base64 form the API
// expects, from whichever of the two content attributes is set.
func resourceTFTPFileData(d *schema.ResourceData) (string, error) {
	if v, ok := d.GetOk("content_base64"); ok {
		raw := v.(string)
		if _, err := base64.StdEncoding.DecodeString(raw); err != nil {
			return "", fmt.Errorf("content_base64 is not valid base64: %w", err)
		}
		return raw, nil
	}
	return base64.StdEncoding.EncodeToString([]byte(d.Get("content").(string))), nil
}

func resourceTFTPFileID(host string, name string) string {
	return fmt.Sprintf("%s/%s", host, name)
}

func resourceTFTPFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	data, err := resourceTFTPFileData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	host := d.Get("host").(string)
	name := d.Get("name").(string)

	hr, err := c.client.RolesTftpAPI.TftpPutFiles(ctx).
		TftpAPIFilesPutInput(api.TftpAPIFilesPutInput{
			Host: host,
			Name: name,
			Data: data,
		}).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	d.SetId(resourceTFTPFileID(host, name))
	return resourceTFTPFileRead(ctx, d, m)
}

func resourceTFTPFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	host := d.Get("host").(string)
	name := d.Get("name").(string)

	// The list endpoint has no filter, so it is only used to confirm the file
	// still exists and to pick up its size.
	res, hr, err := c.client.RolesTftpAPI.TftpGetFiles(ctx).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	var found *api.TftpAPIFile
	for i, f := range res.Files {
		if f.Host == host && f.Name == name {
			found = &res.Files[i]
			break
		}
	}
	if found == nil {
		d.SetId("")
		return diags
	}

	dl, hr, err := c.client.RolesTftpAPI.TftpDownloadFiles(ctx).Host(host).Name(name).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(resourceTFTPFileID(host, name))
	setWrapper(d, "host", found.Host)
	setWrapper(d, "name", found.Name)
	setWrapper(d, "size_bytes", found.SizeBytes)
	// Write back whichever attribute the configuration used, so the contents
	// round-trip in the same encoding they were given in.
	if _, ok := d.GetOk("content_base64"); ok {
		setWrapper(d, "content_base64", dl.Data)
	} else {
		raw, err := base64.StdEncoding.DecodeString(dl.Data)
		if err != nil {
			return diag.FromErr(err)
		}
		setWrapper(d, "content", string(raw))
	}
	return diags
}

func resourceTFTPFileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceTFTPFileCreate(ctx, d, m)
}

func resourceTFTPFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.RolesTftpAPI.TftpDeleteFiles(ctx).
		Host(d.Get("host").(string)).
		Name(d.Get("name").(string)).
		Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}

// resourceTFTPFileImport accepts "<host>/<name>". The name may itself contain
// slashes, so only the first separator is significant.
func resourceTFTPFileImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	host, name, ok := strings.Cut(d.Id(), "/")
	if !ok {
		return nil, fmt.Errorf("invalid ID %q, expected \"<host>/<name>\"", d.Id())
	}
	setWrapper(d, "host", host)
	setWrapper(d, "name", name)
	return []*schema.ResourceData{d}, nil
}
