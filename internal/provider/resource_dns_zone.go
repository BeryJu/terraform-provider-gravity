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
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateMustBeLowercase("DNS name must be lowercase"),
			},
			"authoritative": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"default_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  86400,
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
	m.HandlerConfigs = tfListMap(d.Get("handlers").([]interface{}))
	m.DefaultTTL = int32(d.Get("default_ttl").(int))
	return &m
}

func resourceDNSZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	req := resourceDNSZoneSchemaToModel(d)
	name := d.Get("name").(string)

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
	setWrapper(d, "name", res.Zones[0].Name)
	setWrapper(d, "authoritative", res.Zones[0].Authoritative)
	setWrapper(d, "handlers", res.Zones[0].HandlerConfigs)
	setWrapper(d, "default_ttl", res.Zones[0].DefaultTTL)
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
