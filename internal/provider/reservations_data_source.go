package provider

import (
	"context"
	"fmt"
	"terraform-provider-azureipam/internal/client"
	"terraform-provider-azureipam/internal/gen/data_sources"


	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = (*reservationsDataSource)(nil)

func NewReservationsDataSource() datasource.DataSource {
	return &reservationsDataSource{}
}

type reservationsDataSource struct {
	client *client.Client
}


func (d *reservationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reservations"
}

func (d *reservationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = data_sources.ReservationsDataSourceSchema(ctx)

}

func (d *reservationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data data_sources.ReservationsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.client.ReservationsApiGet(ctx, &data,)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *reservationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Azure IPAM.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
