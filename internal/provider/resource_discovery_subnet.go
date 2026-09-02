package provider

import (
	"context"

	"beryju.io/gravity/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDiscoverySubnet() *schema.Resource {
	return &schema.Resource{
		Description:   "Subnet the discovery role scans for devices.",
		CreateContext: resourceDiscoverySubnetCreate,
		ReadContext:   resourceDiscoverySubnetRead,
		UpdateContext: resourceDiscoverySubnetUpdate,
		DeleteContext: resourceDiscoverySubnetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_resolver": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resolver used to look up hostnames for discovered devices.",
			},
			"discovery_ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     86400,
				Description: "How long a discovered device is remembered for, in seconds.",
			},
		},
	}
}

func resourceDiscoverySubnetSchemaToModel(d *schema.ResourceData) *api.DiscoveryAPISubnetsPutInput {
	return &api.DiscoveryAPISubnetsPutInput{
		SubnetCidr:   d.Get("subnet_cidr").(string),
		DnsResolver:  d.Get("dns_resolver").(string),
		DiscoveryTTL: int32(d.Get("discovery_ttl").(int)),
	}
}

func resourceDiscoverySubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	req := resourceDiscoverySubnetSchemaToModel(d)
	name := d.Get("name").(string)

	hr, err := c.client.RolesDiscoveryAPI.DiscoveryPutSubnets(ctx).
		Identifier(name).
		DiscoveryAPISubnetsPutInput(*req).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	d.SetId(name)
	return resourceDiscoverySubnetRead(ctx, d, m)
}

func resourceDiscoverySubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.RolesDiscoveryAPI.DiscoveryGetSubnets(ctx).Name(d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Subnets) < 1 {
		d.SetId("")
		return diag.Diagnostics{}
	}
	sn := res.Subnets[0]
	setWrapper(d, "name", sn.Name)
	setWrapper(d, "subnet_cidr", sn.SubnetCidr)
	setWrapper(d, "dns_resolver", sn.DnsResolver)
	setWrapper(d, "discovery_ttl", sn.DiscoveryTTL)
	return diags
}

func resourceDiscoverySubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceDiscoverySubnetCreate(ctx, d, m)
}

func resourceDiscoverySubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.RolesDiscoveryAPI.DiscoveryDeleteSubnets(ctx).Identifier(d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
