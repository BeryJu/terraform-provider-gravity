package provider

import (
	"context"

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
			"hook": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Script executed on DHCP events for this scope.",
			},
			"statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"usable": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
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
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Option value as a list of base64-encoded strings.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"value_hex": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Option value as a list of hex-encoded strings.",
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
		Ttl:        int64(d.Get("lease_ttl").(int)),
		Options:    []api.TypesDHCPOption{},
	}
	m.Hook = d.Get("hook").(string)
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
			aopt.Value64 = sliceToString(t)
		}
		if t, ok := values["value_hex"].([]interface{}); ok && len(t) > 0 {
			aopt.ValueHex = sliceToString(t)
		}
		m.Options = append(m.Options, aopt)
	}

	// Only send a DNS object when the block is actually configured: the server
	// stores and returns exactly what it is given, so sending an empty object
	// for an absent block would show up as drift on the next refresh.
	dns := d.Get("dns").(*schema.Set)
	for _, opt := range dns.List() {
		values := opt.(map[string]interface{})
		m.Dns = &api.DhcpScopeDNS{}

		if t, ok := values["add_zone_in_hostname"].(bool); ok {
			m.Dns.AddZoneInHostname = api.PtrBool(t)
		}
		if t, ok := values["zone"].(string); ok && t != "" {
			m.Dns.Zone = api.PtrString(t)
		}
		if t, ok := values["search"].([]interface{}); ok && len(t) > 0 {
			m.Dns.Search = sliceToString(t)
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

// flattenOptions and flattenDNS return plain []interface{} rather than a
// *schema.Set. Building the set here needs a hash function that matches the
// one the SDK derives from the schema; supplying our own instead silently
// collapsed options that differed only in fields the hash failed to read,
// leaving a permanent diff. Handing the SDK a slice lets it hash correctly.
func flattenOptions(opts []api.TypesDHCPOption) []interface{} {
	vopts := make([]interface{}, len(opts))
	for i, opt := range opts {
		vopts[i] = map[string]interface{}{
			"tag":       int32Value(opt.Tag.Get()),
			"tag_name":  stringValue(opt.TagName),
			"value":     stringValue(opt.Value.Get()),
			"value64":   stringSlice(opt.Value64),
			"value_hex": stringSlice(opt.ValueHex),
		}
	}
	return vopts
}

func flattenDNS(dns *api.DhcpScopeDNS) []interface{} {
	if dns == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"add_zone_in_hostname": boolValue(dns.AddZoneInHostname),
			"zone":                 stringValue(dns.Zone),
			"search":               stringSlice(dns.Search),
		},
	}
}

func flattenScopeStatistics(st api.DhcpAPIScopeStatistics) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"usable": int(st.Usable),
			"used":   int(st.Used),
		},
	}
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
	sc := res.Scopes[0]
	setWrapper(d, "name", sc.Scope)
	setWrapper(d, "default", sc.Default)
	setWrapper(d, "subnet_cidr", sc.SubnetCidr)
	setWrapper(d, "lease_ttl", sc.Ttl)
	setWrapper(d, "hook", sc.Hook)
	setWrapper(d, "ipam", sc.Ipam)
	setWrapper(d, "option", flattenOptions(sc.Options))
	setWrapper(d, "dns", flattenDNS(sc.Dns))
	setWrapper(d, "statistics", flattenScopeStatistics(sc.Statistics))
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
