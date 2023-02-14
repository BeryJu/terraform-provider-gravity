package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"beryju.io/gravity/api"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Computed {
			desc += " Generated."
		}
		return strings.TrimSpace(desc)
	}
}

// Provider -
func Provider(version string, testing bool) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAVITY_URL", nil),
				Description: "The Gravity API endpoint, can optionally be passed as `GRAVITY_URL` environmental variable",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAVITY_INSECURE", false),
				Description: "Whether to skip TLS verification, can optionally be passed as `GRAVITY_INSECURE` environmental variable",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAVITY_TOKEN", nil),
				Sensitive:   true,
				Description: "The Gravity API token, can optionally be passed as `GRAVITY_TOKEN` environmental variable",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gravity_dns_zone":   resourceDNSZone(),
			"gravity_dns_record": resourceDNSRecord(),
			"gravity_dhcp_scope": resourceDHCPScope(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure(version, testing),
	}
}

// APIClient Hold the API Client and any relevant configuration
type APIClient struct {
	client *api.APIClient
}

func providerConfigure(version string, testing bool) schema.ConfigureContextFunc {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		apiURL := d.Get("url").(string)
		token := d.Get("token").(string)
		insecure := d.Get("insecure").(bool)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		akURL, err := url.Parse(apiURL)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		config := api.NewConfiguration()
		config.Debug = true
		config.UserAgent = fmt.Sprintf("gravity-terraform@%s", version)
		config.Host = akURL.Host
		config.Scheme = akURL.Scheme
		if testing {
			config.HTTPClient = &http.Client{
				Transport: NewTestingTransport(GetTLSTransport(insecure)),
			}
		} else {
			config.HTTPClient = &http.Client{
				Transport: GetTLSTransport(insecure),
			}
		}

		config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
		apiClient := api.NewAPIClient(config)

		return &APIClient{
			client: apiClient,
		}, diags
	}
}

// TestingTransport Transport used for testing, always returns a 400 Response
type TestingTransport struct {
	inner http.RoundTripper
}

// NewTestingTransport Get a HTTP Transport that fails all requests
func NewTestingTransport(inner http.RoundTripper) *TestingTransport {
	return &TestingTransport{inner}
}

// RoundTrip HTTP Transport
func (tt *TestingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	body := "mock-failed-request"
	return &http.Response{
		Status:        "400 Bad Request",
		StatusCode:    400,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       r,
		Header:        make(http.Header),
	}, nil
}

// GetTLSTransport Get a TLS transport instance, that skips verification if configured via environment variables.
func GetTLSTransport(insecure bool) http.RoundTripper {
	tlsTransport, err := httptransport.TLSTransport(httptransport.TLSClientOptions{
		InsecureSkipVerify: insecure,
	})
	if err != nil {
		panic(err)
	}
	return tlsTransport
}
