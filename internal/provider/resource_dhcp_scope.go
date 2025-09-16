package provider

import (
	"bytes"
	"context"
	"fmt"

	"beryju.io/gravity/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDHCPScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDHCPScopeCreate,
		ReadContext:   resourceDHCPScopeRead,
		UpdateContext: resourceDHCPScopeUpdate,
		DeleteContext: resourceDHCPScopeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"subnet_cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lease_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  86400,
			},
			"ipam": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dns": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"add_zone_in_hostname": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"search": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"option": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"tag_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value64": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceDHCPScopeSchemaToModel(d *schema.ResourceData) (*api.DhcpAPIScopesPutInput, error) {
	m := api.DhcpAPIScopesPutInput{
		Default:    d.Get("default").(bool),
		SubnetCidr: d.Get("subnet_cidr").(string),
		Ttl:        int32(d.Get("lease_ttl").(int)),
		Options:    []api.TypesDHCPOption{},
		Dns:        &api.DhcpScopeDNS{},
	}
	m.Ipam = tfMap(d.Get("ipam").(map[string]interface{}))

	options := d.Get("option").(*schema.Set)
	for _, opt := range options.List() {
		values := opt.(map[string]interface{})
		aopt := api.TypesDHCPOption{}

		if t, ok := values["tag"].(int); ok && t > 0 {
			aopt.Tag.Set(api.PtrInt32(int32(t)))
		}
		if t, ok := values["tag_name"].(string); ok && t != "" {
			aopt.TagName = api.PtrString(t)
		}
		if t, ok := values["value"].(string); ok && t != "" {
			aopt.Value.Set(api.PtrString(t))
		}
		if t, ok := values["value64"].([]interface{}); ok && len(t) > 0 {
			values := make([]string, len(t))
			for i, v := range t {
				values[i] = v.(string)
			}
			aopt.Value64 = values
		}
		m.Options = append(m.Options, aopt)
	}

	dns := d.Get("dns").(*schema.Set)
	for _, opt := range dns.List() {
		values := opt.(map[string]interface{})

		if t, ok := values["add_zone_in_hostname"].(bool); ok {
			m.Dns.AddZoneInHostname = api.PtrBool(t)
		}
		if t, ok := values["zone"].(string); ok && t != "" {
			m.Dns.Zone = api.PtrString(t)
		}
		if t, ok := values["search"].([]interface{}); ok && len(t) > 0 {
			values := make([]string, len(t))
			for i, v := range t {
				values[i] = v.(string)
			}
			m.Dns.Search = values
		}
	}
	return &m, nil
}

func resourceDHCPScopeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	req, err := resourceDHCPScopeSchemaToModel(d)
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)

	hr, err := c.client.RolesDhcpAPI.DhcpPutScopes(ctx).Scope(name).DhcpAPIScopesPutInput(*req).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	d.SetId(name)
	return resourceDHCPScopeRead(ctx, d, m)
}

func flattenOptions(opts []api.TypesDHCPOption) *schema.Set {
	var vopts []interface{}

	for _, opt := range opts {
		vopt := map[string]interface{}{}
		vopt["tag"] = opt.Tag.Get()
		vopt["tag_name"] = opt.TagName
		vopt["value"] = opt.Value.Get()
		if len(opt.Value64) > 0 {
			vopt["value64"] = opt.Value64
		}
		vopts = append(vopts, vopt)
	}

	return schema.NewSet(func(i interface{}) int {
		var buf bytes.Buffer
		mCondition := i.(map[string]interface{})
		if v, ok := mCondition["tag"].(int); ok {
			buf.WriteString(fmt.Sprintf("%d-", v))
		}
		if v, ok := mCondition["tag_name"].(string); ok {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
		if v, ok := mCondition["value"].(string); ok {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
		if v, ok := mCondition["value64"].([]string); ok {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
		return StringHashcode(buf.String())
	}, vopts)
}

func flattenDNS(dns *api.DhcpScopeDNS) *schema.Set {
	var vdns []interface{}

	if dns != nil {
		vopt := map[string]interface{}{}
		vopt["add_zone_in_hostname"] = dns.AddZoneInHostname
		vopt["zone"] = dns.Zone
		if len(dns.Search) > 0 {
			vopt["search"] = dns.Search
		}
		vdns = append(vdns, vopt)
	}

	return schema.NewSet(func(i interface{}) int {
		var buf bytes.Buffer
		mCondition := i.(map[string]interface{})
		if v, ok := mCondition["add_zone_in_hostname"].(bool); ok {
			buf.WriteString(fmt.Sprintf("%t-", v))
		}
		if v, ok := mCondition["zone"].(string); ok {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
		if v, ok := mCondition["search"].([]string); ok {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
		return StringHashcode(buf.String())
	}, vdns)
}

func resourceDHCPScopeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.RolesDhcpAPI.DhcpGetScopes(ctx).Name(d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Scopes) < 1 {
		d.SetId("")
		return diag.Diagnostics{}
	}
	setWrapper(d, "name", res.Scopes[0].Scope)
	setWrapper(d, "default", res.Scopes[0].Default)
	setWrapper(d, "subnet_cidr", res.Scopes[0].SubnetCidr)
	setWrapper(d, "lease_ttl", res.Scopes[0].Ttl)
	setWrapper(d, "ipam", res.Scopes[0].Ipam)
	setWrapper(d, "option", flattenOptions(res.Scopes[0].Options))
	setWrapper(d, "dns", flattenDNS(res.Scopes[0].Dns))
	return diags
}

func resourceDHCPScopeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := resourceDHCPScopeCreate(ctx, d, m)
	if diag != nil {
		return diag
	}
	return resourceDHCPScopeRead(ctx, d, m)
}

func resourceDHCPScopeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.RolesDhcpAPI.DhcpDeleteScopes(ctx).Scope(d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
