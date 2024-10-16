package provider

import (
	"context"

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
			StateContext: schema.ImportStatePassthroughContext,
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
		},
	}
}

func resourceDHCPLeaseSchemaToModel(d *schema.ResourceData) *api.DhcpAPILeasesPutInput {
	m := api.DhcpAPILeasesPutInput{
		Address:  d.Get("address").(string),
		Hostname: d.Get("hostname").(string),
	}
	if res := d.Get("reservation").(bool); res {
		m.Expiry = api.PtrInt32(-1)
	}
	return &m
}

func resourceDHCPLeaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	req := resourceDHCPLeaseSchemaToModel(d)
	scope := d.Get("scope").(string)
	identifier := d.Get("identifier").(string)

	hr, err := c.client.RolesDhcpApi.DhcpPutLeases(ctx).
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

	res, hr, err := c.client.RolesDhcpApi.DhcpGetLeases(ctx).
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
	d.SetId(res.Leases[0].Identifier)
	setWrapper(d, "identifier", res.Leases[0].Identifier)
	setWrapper(d, "address", res.Leases[0].Address)
	setWrapper(d, "hostname", res.Leases[0].Hostname)
	setWrapper(d, "scope", res.Leases[0].ScopeKey)
	setWrapper(d, "reservation", *res.Leases[0].Expiry <= -1)
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

	hr, err := c.client.RolesDhcpApi.DhcpDeleteLeases(ctx).
		Scope(scope).
		Identifier(identifier).
		Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
