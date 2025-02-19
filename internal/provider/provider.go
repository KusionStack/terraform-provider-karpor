package provider

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider                   = &KarporProvider{}
	_ provider.ProviderWithValidateConfig = &KarporProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KarporProvider{
			version: version,
		}
	}
}

// KarporProvider is the provider implementation.
type KarporProvider struct {
	version string
}

// Metadata returns the provider type name.
func (p *KarporProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "karpor"
	resp.Version = p.version
}

// Schema returns the provider schema.
func (p *KarporProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "Karpor API endpoint URL",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "API key for authentication",
			},
			"skip_tls_verify": schema.BoolAttribute{
				Optional:    true,
				Description: "Skip TLS verification, by default it is false",
			},
		},
	}
}

// Configure configures the provider.
func (p *KarporProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config KarporProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	api_endpoint := os.Getenv("KARPOR_API_ENDPOINT")
	api_key := os.Getenv("KARPOR_API_KEY")
	skip_tls_verify := false
	if strings.ToLower(os.Getenv("KARPOR_SKIP_TLS_VERIFY")) == "true" {
		skip_tls_verify = true
	}

	if !config.ApiEndpoint.IsNull() {
		api_endpoint = config.ApiEndpoint.ValueString()
	}
	if !config.ApiKey.IsNull() {
		api_key = config.ApiKey.ValueString()
	}
	if !config.SkipTlsVerify.IsNull() {
		skip_tls_verify = config.SkipTlsVerify.ValueBool()
	}

	if api_endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_endpoint"),
			"Missing Karpor API Endpoint",
			"The provider cannot create the Karpor API client as there is a missing or empty configuration value for the Karpor API endpoint.",
		)
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Karpor API Key",
			"The provider cannot create the Karpor API client as there is a missing or empty configuration value for the Karpor API key.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "endpoint", api_endpoint)
	ctx = tflog.SetField(ctx, "key", api_key)
	ctx = tflog.SetField(ctx, "skip_tls_verify", skip_tls_verify)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "key")

	tflog.Debug(ctx, "Creating Karpor client")

	client, err := NewKarporClient(api_endpoint, api_key, skip_tls_verify)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Karpor client",
			"An unexpected error occurred when creating the Karpor client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Karpor Client Error: "+err.Error(),
		)
		return
	}

	// Make client available during data source and resource operations
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Debug(ctx, "Karpor client created successfully")
}

// Resources returns the resources supported by the provider.
func (p *KarporProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterRegistrationResource,
	}
}

// DataSources returns the data sources supported by the provider.
func (p *KarporProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClusterDataSource,
	}
}

// ValidateConfig validates the provider configuration.
func (p *KarporProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var config KarporProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiEndpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_endpoint"),
			"Unknown Karpor API Endpoint",
			"The provider cannot create the Karpor API client as there is an unknown configuration value for the Karpor API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the KARPOR_API_ENDPOINT environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Karpor API Key",
			"The provider cannot create the Karpor API client as there is an unknown configuration value for the Karpor API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_PASSWORD environment variable.",
		)
	}

	if config.SkipTlsVerify.ValueBool() {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("skip_tls_verify"),
			"Karpor Skip TLS Verify",
			"Karpor Skip TLS Verify is not recommended for production environments. It is recommended to set the value to true only for development purposes.",
		)
	}
}

// ProviderModel is the provider model.
type KarporProviderModel struct {
	ApiEndpoint   types.String `tfsdk:"api_endpoint"`
	ApiKey        types.String `tfsdk:"api_key"`
	SkipTlsVerify types.Bool   `tfsdk:"skip_tls_verify"`
}
