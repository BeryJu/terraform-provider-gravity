package provider

import (
	"context"

	"beryju.io/gravity/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSZoneCreate,
		ReadContext:   resourceDNSZoneRead,
		UpdateContext: resourceDNSZoneUpdate,
		DeleteContext: resourceDNSZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authoritative": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"handlers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func resourceDNSZoneSchemaToModel(d *schema.ResourceData) *api.DnsAPIZonesPutInput {
	m := api.DnsAPIZonesPutInput{}
	m.Authoritative = d.Get("authoritative").(bool)

	rawHandlers := d.Get("handlers").([]interface{})
	handlers := make([]map[string]string, len(rawHandlers))
	for i, rh := range rawHandlers {
		m := make(map[string]string)
		for k, v := range rh.(map[string]interface{}) {
			m[k] = v.(string)
		}
		handlers[i] = m
	}
	m.HandlerConfigs = handlers
	return &m
}

func resourceDNSZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	req := resourceDNSZoneSchemaToModel(d)
	name := d.Get("zone").(string)

	hr, err := c.client.RolesDnsApi.DnsPutZones(ctx).Zone(name).DnsAPIZonesPutInput(*req).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	d.SetId(name)
	return resourceDNSZoneRead(ctx, d, m)
}

func resourceDNSZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.RolesDnsApi.DnsGetZones(ctx).Name(d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Zones) < 1 {
		d.SetId("")
		return diag.Diagnostics{}
	}
	setWrapper(d, "zone", res.Zones[0].Name)
	setWrapper(d, "authoritative", res.Zones[0].Authoritative)
	setWrapper(d, "handlers", res.Zones[0].HandlerConfigs)
	return diags
}

func resourceDNSZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := resourceDNSZoneCreate(ctx, d, m)
	if diag != nil {
		return diag
	}
	return resourceDNSZoneRead(ctx, d, m)
}

func resourceDNSZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.RolesDnsApi.DnsDeleteZones(ctx).Zone(d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
