package provider

import (
	"context"
	"fmt"
	"math/big"
	"terraform-provider-azureipam/internal/client"
	"terraform-provider-azureipam/internal/gen/data_sources"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	resp.Diagnostics.Append(reservationsDataGet(ctx, d.client, &data)...)
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

func reservationsDataGet(ctx context.Context, c *client.Client, data *data_sources.ReservationsModel) diag.Diagnostics {
	var diags diag.Diagnostics
	requestData := client.ReservationApiModel{
		Space: data.Space.ValueString(),
		Block: data.Block.ValueString(),
		Settled: data.Settled.ValueBool(),
	}
	response, err := c.ReservationsApiGet(ctx, requestData)
	if err != nil {
		diags.AddError("Response Unmarshal Error", fmt.Sprintf("Failed to unmarshal response: %s", err))
		return diags
	}
	elements := make([]attr.Value, len(response))
	for i, reservation := range response {
		tags := make(map[string]attr.Value, len(reservation.Tag))
		for k, v := range reservation.Tag {
			tags[k] = types.StringValue(v)
		}
		createdOn := big.NewFloat(reservation.CreatedOn)
		settledOn := big.NewFloat(reservation.SettledOn)
		tag, _ := types.MapValue(types.StringType, tags)
		objVal := data_sources.NewReservationsValueMust(data_sources.NewReservationsValueNull().AttributeTypes(ctx),
			map[string]attr.Value{
				"id":         types.StringValue(reservation.Id),
				"space":      types.StringValue(reservation.Space),
				"block":      types.StringValue(reservation.Block),
				"cidr":       types.StringValue(reservation.CIDR),
				"desc":       types.StringValue(reservation.Desc),
				"created_on": types.NumberValue(createdOn),
				"created_by": types.StringValue(reservation.CreatedBy),
				"settled_by": types.StringValue(reservation.SettledBy),
				"settled_on": types.NumberValue(settledOn),
				"status":     types.StringValue(reservation.Status),
				"tag":        tag,
			},
		)
		// diags.Append(objDiags...)
		// if diags.HasError() {
		// 	return diags
		// }
		elements[i] = objVal
	}

	// Set the Reservations field in the ReservationsModel
	data.Reservations, diags = types.SetValueFrom(ctx, data_sources.NewReservationsValueNull().Type(ctx), &elements)

	return diags

}
