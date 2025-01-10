package provider

import (
	"context"
	"fmt"

	"beryju.io/gravity/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecordCreate,
		ReadContext:   resourceDNSRecordRead,
		UpdateContext: resourceDNSRecordUpdate,
		DeleteContext: resourceDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateMustBeLowercase("DNS name must be lowercase"),
			},
			"hostname": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateMustBeLowercase("DNS name must be lowercase"),
			},
			"uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"mx_preference": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"srv_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"srv_priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"srv_weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDNSRecordSchemaToModel(d *schema.ResourceData) *api.DnsAPIRecordsPutInput {
	m := api.DnsAPIRecordsPutInput{
		Data: d.Get("data").(string),
		Type: d.Get("type").(string),
	}
	if v, ok := d.GetOk("mx_preference"); ok {
		va := v.(int)
		m.MxPreference = api.PtrInt32(int32(va))
	}
	if v, ok := d.GetOk("srv_port"); ok {
		va := v.(int)
		m.SrvPort = api.PtrInt32(int32(va))
	}
	if v, ok := d.GetOk("srv_priority"); ok {
		va := v.(int)
		m.SrvPriority = api.PtrInt32(int32(va))
	}
	if v, ok := d.GetOk("srv_weight"); ok {
		va := v.(int)
		m.SrvWeight = api.PtrInt32(int32(va))
	}
	return &m
}

func resourceDNSRecordID(d *schema.ResourceData) string {
	zone := d.Get("zone").(string)
	hostname := d.Get("hostname").(string)
	type_ := d.Get("type").(string)
	uid := d.Get("uid").(string)
	return fmt.Sprintf("%s:%s:%s:%s", zone, hostname, type_, uid)
}

func resourceDNSRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	req := resourceDNSRecordSchemaToModel(d)
	zone := d.Get("zone").(string)
	hostname := d.Get("hostname").(string)
	uid := d.Get("uid").(string)

	hr, err := c.client.RolesDnsApi.DnsPutRecords(ctx).
		Zone(zone).
		Hostname(hostname).
		Uid(uid).
		DnsAPIRecordsPutInput(*req).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	d.SetId(resourceDNSRecordID(d))
	return resourceDNSRecordRead(ctx, d, m)
}

func resourceDNSRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	zone := d.Get("zone").(string)
	hostname := d.Get("hostname").(string)
	type_ := d.Get("type").(string)
	uid := d.Get("uid").(string)

	res, hr, err := c.client.RolesDnsApi.DnsGetRecords(ctx).
		Zone(zone).
		Hostname(hostname).
		Type_(type_).
		Uid(uid).
		Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Records) < 1 {
		d.SetId("")
		return diag.Diagnostics{}
	}
	d.SetId(resourceDNSRecordID(d))
	setWrapper(d, "fqdn", res.Records[0].Fqdn)
	setWrapper(d, "uid", res.Records[0].Uid)
	setWrapper(d, "hostname", res.Records[0].Hostname)
	setWrapper(d, "data", res.Records[0].Data)
	setWrapper(d, "type", res.Records[0].Type)
	setWrapper(d, "mx_preference", res.Records[0].MxPreference)
	setWrapper(d, "srv_port", res.Records[0].SrvPort)
	setWrapper(d, "srv_priority", res.Records[0].SrvPriority)
	setWrapper(d, "srv_weight", res.Records[0].SrvWeight)
	return diags
}

func resourceDNSRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diag := resourceDNSRecordCreate(ctx, d, m)
	if diag != nil {
		return diag
	}
	return resourceDNSRecordRead(ctx, d, m)
}

func resourceDNSRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	zone := d.Get("zone").(string)
	hostname := d.Get("hostname").(string)
	uid := d.Get("uid").(string)
	type_ := d.Get("type").(string)

	hr, err := c.client.RolesDnsApi.DnsDeleteRecords(ctx).
		Zone(zone).
		Hostname(hostname).
		Uid(uid).
		Type_(type_).
		Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
