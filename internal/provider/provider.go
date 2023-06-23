package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/genproto/googleapis/datastore/v1"

	"github.com/Mirantis/terraform-provider-mke/internal/client"
)

const (
	MKEProviderName = "mke"
)

// Ensure MKEProvider satisfies various provider interfaces.
var _ provider.Provider = &MKEProvider{}

// MKEProvider defines the provider implementation.
type MKEProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MKEProviderModel describes the provider data model.
type MKEProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	UnsafeSSL types.Bool `tfsdk:"unsafe_ssl_client"`
}

// Client MKE client generation
func (pm MKEProviderModel) Client() (client.Client, error) {
	if pm.UnsafeSSL {
		return client.NewUnsafeSSLClient(pm.Endpoint.ValueString(), pm.Username.ValueString(), pm.Password.ValueString())
	}
	return client.NewClientSimple(pm.Endpoint.ValueString(), pm.Username.ValueString(), pm.Password.ValueString())
}

// Metadata Terraform metadata handler
func (p *MKEProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = MKEProvider
	resp.Version = p.version
}

func (p *MKEProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "MKE Endpoint for API access.",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for MKE API access.",
				Required: true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "User password for MKE API access.",
				Required: true,
				Sensitive: true,
			},
			"unsafe_ssl_client": schema.BoolAttribute{
				MarkdownDescription: "Allow connections without SSL validation on the http connection.",
				Optional: true,
			},
		},
	}
}

func (p *MKEProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MKEProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = &data
	resp.ResourceData = &data
}

func (p *MKEProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		ClientBundleResource,
	}
}

func (p *MKEProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MKEProvider{
			version: version,
		}
	}
}
