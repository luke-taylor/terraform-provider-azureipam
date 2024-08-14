package provider

import (
	"context"
	"os"
	"terraform-provider-azureipam/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type azureipamProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Token    types.String `tfsdk:"token"`
	ClientID types.String `tfsdk:"client_id"`
}

func (p *azureipamProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with IPAM.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "uri for ipam azure.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "bearer token for azure ipam.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "Client Id for the Azure IPAM Engine App.",
				Optional:    true,
			},
		},
	}
}

func (p *azureipamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring HashiCups client")
	var config azureipamProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown HashiCups API Host",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown HashiCups API token",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_USERNAME environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("HASHICUPS_HOST")
	token := os.Getenv("HASHICUPS_TOKEN")
	clientID := os.Getenv("HASHICUPS_CLIENT_ID")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if !config.ClientID.IsNull() {
		clientID = config.ClientID.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing HashiCups API Host",
			"The provider cannot create the HashiCups API client as there is a missing or empty value for the HashiCups API host. "+
				"Set the host value in the configuration or use the HASHICUPS_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	azToken, _ := getAzureAccessToken(clientID)
	if token == "" {
		token = azToken
	}
	// Create a new HashiCups client using the configuration values
	client, err := client.NewClient(&host, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create HashiCups API Client",
			"An unexpected error occurred when creating the HashiCups API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"HashiCups Client Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
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
	return []func() resource.Resource{}
}
