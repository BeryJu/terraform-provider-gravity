package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"beryju.io/gravity/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure GravityProvider satisfies various provider interfaces.
var _ provider.Provider = &GravityProvider{}

// GravityProvider defines the provider implementation.
type GravityProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// GravityProviderModel describes the provider data model.
type GravityProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func (p *GravityProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gravity"
	resp.Version = p.version
}

func (p *GravityProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Gravity API Endpoint, listens by default on port 9009.",
			},
			"insecure": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Set to true to ignore any HTTPS Certificate errors",
			},
			"token": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Token to authenticate against the API with.",
			},
		},
	}
}

func (p *GravityProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Check environment variables
	apiToken := os.Getenv("GRAVITY_TOKEN")
	endpoint := os.Getenv("GRAVITY_ENDPOINT")

	var data GravityProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Token.ValueString() != "" {
		apiToken = data.Token.ValueString()
	}

	if data.Endpoint.ValueString() != "" {
		endpoint = data.Endpoint.ValueString()
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		resp.Diagnostics.AddError("failed to parse endpoint", err.Error())
		return
	}

	config := api.NewConfiguration()
	config.Debug = true
	config.UserAgent = fmt.Sprintf("terraform-provider-gravity@%s", p.version)
	config.Host = url.Host
	config.Scheme = url.Scheme
	config.HTTPClient = &http.Client{
		Transport: GetTLSTransport(data.Insecure.ValueBool()),
	}

	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	apiClient := api.NewAPIClient(config)

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *GravityProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDNSZoneResource,
	}
}

func (p *GravityProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GravityProvider{
			version: version,
		}
	}
}
