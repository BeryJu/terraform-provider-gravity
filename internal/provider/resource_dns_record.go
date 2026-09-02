package provider

import (
	"context"
	"fmt"
	"strings"

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
			StateContext: resourceDNSRecordImport,
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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      EnumToDescription(api.AllowedTypesDNSRecordTypeEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedTypesDNSRecordTypeEnumValues),
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
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

			"soa_mbox": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"soa_serial": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"soa_refresh": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"soa_retry": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"soa_expire": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDNSRecordSchemaToModel(d *schema.ResourceData) *api.DnsAPIRecordsPutInput {
	m := api.DnsAPIRecordsPutInput{
		Data: d.Get("data").(string),
		Type: api.TypesDNSRecordType(d.Get("type").(string)),
		Ttl:  int64(d.Get("ttl").(int)),
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
	if v, ok := d.GetOk("soa_mbox"); ok {
		m.SoaMbox = api.PtrString(v.(string))
	}
	if v, ok := d.GetOk("soa_serial"); ok {
		m.SoaSerial = api.PtrInt32(int32(v.(int)))
	}
	if v, ok := d.GetOk("soa_refresh"); ok {
		m.SoaRefresh = api.PtrInt32(int32(v.(int)))
	}
	if v, ok := d.GetOk("soa_retry"); ok {
		m.SoaRetry = api.PtrInt32(int32(v.(int)))
	}
	if v, ok := d.GetOk("soa_expire"); ok {
		m.SoaExpire = api.PtrInt32(int32(v.(int)))
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

	hr, err := c.client.RolesDnsAPI.DnsPutRecords(ctx).
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

	res, hr, err := c.client.RolesDnsAPI.DnsGetRecords(ctx).
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
	r := res.Records[0]
	setWrapper(d, "zone", zone)
	setWrapper(d, "fqdn", r.Fqdn)
	setWrapper(d, "uid", r.Uid)
	setWrapper(d, "hostname", r.Hostname)
	setWrapper(d, "data", r.Data)
	setWrapper(d, "type", string(r.Type))
	setWrapper(d, "ttl", r.Ttl)
	setWrapper(d, "mx_preference", int32Value(r.MxPreference))
	setWrapper(d, "srv_port", int32Value(r.SrvPort))
	setWrapper(d, "srv_priority", int32Value(r.SrvPriority))
	setWrapper(d, "srv_weight", int32Value(r.SrvWeight))
	setWrapper(d, "soa_mbox", stringValue(r.SoaMbox))
	setWrapper(d, "soa_serial", int32Value(r.SoaSerial))
	setWrapper(d, "soa_refresh", int32Value(r.SoaRefresh))
	setWrapper(d, "soa_retry", int32Value(r.SoaRetry))
	setWrapper(d, "soa_expire", int32Value(r.SoaExpire))
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
	type_ := api.TypesDNSRecordType(d.Get("type").(string))

	hr, err := c.client.RolesDnsAPI.DnsDeleteRecords(ctx).
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

// resourceDNSRecordImport splits the composite ID back into the four
// attributes Read needs, since none of them are derivable from the API alone.
func resourceDNSRecordImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid ID %q, expected \"<zone>:<hostname>:<type>:<uid>\"", d.Id())
	}
	setWrapper(d, "zone", parts[0])
	setWrapper(d, "hostname", parts[1])
	setWrapper(d, "type", parts[2])
	setWrapper(d, "uid", parts[3])
	return []*schema.ResourceData{d}, nil
}
