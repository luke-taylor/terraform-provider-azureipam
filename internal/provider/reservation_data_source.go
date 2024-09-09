package provider

import (
	"context"
	"fmt"
	"terraform-provider-azureipam/internal/client"
	"terraform-provider-azureipam/internal/gen/data_sources"

	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*reservationDataSource)(nil)

func NewReservationDataSource() datasource.DataSource {
	return &reservationDataSource{}
}

type reservationDataSource struct {
	client *client.Client
}

func (d *reservationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reservation"
}

func (d *reservationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = data_sources.ReservationDataSourceSchema(ctx)
}

func (d *reservationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data data_sources.ReservationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(reservationApiGet(ctx, d.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
func (d *reservationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func reservationApiGet(ctx context.Context, c *client.Client, data *data_sources.ReservationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	requestData := client.ReservationApiModel{
		Space: data.Space.ValueString(),
		Block: data.Block.ValueString(),
		Id:    data.Id.ValueString(),
	}

	response, err := c.ReservationApiGet(ctx, requestData)
	if err != nil {
		diags.AddError("Response Unmarshal Error", fmt.Sprintf("Failed to unmarshal response: %s", err))
		return diags
	}

	data.Id = types.StringValue(response.Id)
	data.Cidr = types.StringValue(response.CIDR)
	data.CreatedBy = types.StringValue(response.CreatedBy)
	data.Desc = types.StringValue(response.Desc)
	data.Status = types.StringValue(response.Status)
	settledOnBigFloat := big.NewFloat(response.SettledOn)
	data.SettledOn = types.NumberValue(settledOnBigFloat)
	data.SettledBy = types.StringValue(response.SettledBy)
	createdOnBigFloat := big.NewFloat(response.CreatedOn)
	data.CreatedOn = types.NumberValue(createdOnBigFloat)

	if response.Tag != nil {
		tagElements := make(map[string]attr.Value)
		for k, v := range response.Tag {
			tagElements[k] = types.StringValue(v)
		}
		data.Tag, _ = types.MapValue(types.StringType, tagElements)
	} else {
		data.Tag = types.MapNull(types.StringType)
	}
	return diags
}
