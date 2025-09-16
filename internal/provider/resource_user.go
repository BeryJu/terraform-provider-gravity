package provider

import (
	"context"
	"encoding/json"
	"strings"

	"beryju.io/gravity/api"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permissions": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "[]",
				ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
					err := json.Unmarshal([]byte(i.(string)), &[]struct{}{})
					if err != nil {
						return diag.FromErr(errors.Wrap(err, "Failed to validate handlers"))
					}
					return nil
				},
			},
		},
	}
}

func resourceUserSchemaToModel(d *schema.ResourceData) (*api.AuthAPIUsersPutInput, diag.Diagnostics) {
	m := api.AuthAPIUsersPutInput{}

	var c []api.AuthPermission
	err := json.NewDecoder(strings.NewReader(d.Get("permissions").(string))).Decode(&c)
	if err != nil {
		return nil, diag.FromErr(errors.Wrap(err, "failed to convert to json"))
	}
	m.Permissions = c

	return &m, nil
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	req, diags := resourceUserSchemaToModel(d)
	if diags != nil {
		return diags
	}
	username := d.Get("username").(string)

	hr, err := c.client.RolesApiAPI.ApiPutUsers(ctx).Username(username).AuthAPIUsersPutInput(*req).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	d.SetId(username)
	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.RolesApiAPI.ApiGetUsers(ctx).Username(d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Users) < 1 {
		d.SetId("")
		return diag.Diagnostics{}
	}
	setWrapper(d, "username", res.Users[0].Username)
	b, err := json.Marshal(res.Users[0].Permissions)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "permissions", string(b))
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := resourceUserCreate(ctx, d, m)
	if diag != nil {
		return diag
	}
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.RolesApiAPI.ApiDeleteUsers(ctx).Username(d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
