package provider

import (
	"context"
	"os"
	"terraform-provider-azureipam/internal/client"
	"terraform-provider-azureipam/internal/gen"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = (*azureipamProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &azureipamProvider{
			version: version,
		}
	}
}

type azureipamProvider struct {
	version string
}

func (p *azureipamProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = gen.AzureipamProviderSchema(ctx)
}

func (p *azureipamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Azure IPAM client")
	var config gen.AzureipamModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.HostUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host_url"),
			"Unknown Azure IPAM Solution URL",
			"The provider cannot create the Azure IPAM client as there is an unknown configuration value for the Host URL "+
				"Set the value statically in the configuration, or use the IPAM_HOST_URL environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Azure IPAM Solution Token",
			"The provider cannot create the Azure IPAM client as there is an unknown configuration value for the Token "+
				"Set the value statically in the configuration, or use the IPAM_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("IPAM_HOST_URL")
	token := os.Getenv("IPAM_TOKEN")
	clientID := os.Getenv("IPAM_ENGINE_CLIENT_ID")

	if !config.HostUrl.IsNull() {
		host = config.HostUrl.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if !config.EngineClientId.IsNull() {
		clientID = config.EngineClientId.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Azure IPAM Solution URL",
			"The provider cannot create the Azure IPAM client as there is a missing configuration value for the Host URL "+
				"Set the value statically in the configuration, or use the IPAM_HOST_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	azToken, _ := client.GetAzureAccessToken(clientID)
	if token == "" {
		token = azToken
	}
	// Create a new Azure IPAM client using the configuration values
	client, err := client.NewClient(&host, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Azure IPAM client",
			"Failed to create Azure IPAM client: "+err.Error(),
		)
		return
	}

	// Make the Azure IPAM client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *azureipamProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "azureipam"
}

func (p *azureipamProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAdminsDataSource,
	}
}

func (p *azureipamProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewReservationResource,
	}
}
