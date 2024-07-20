package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*azureipamProvider)(nil)

func New() func() provider.Provider {
	return func() provider.Provider {
		return &azureipamProvider{}
	}
}

type azureipamProvider struct{}

func (p *azureipamProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {

}

func (p *azureipamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *azureipamProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "azureipam"
}

func (p *petstoreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
			NewAdminsDataSource,
	}
}

func (p *azureipamProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
