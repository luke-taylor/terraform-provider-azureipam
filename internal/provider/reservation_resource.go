package provider

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"terraform-provider-azureipam/internal/client"
	"terraform-provider-azureipam/internal/gen/resources"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*reservationResource)(nil)

func NewReservationResource() resource.Resource {
	return &reservationResource{}
}

type reservationResource struct {
	client *client.Client
}

func (r *reservationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reservation"
}

func (r *reservationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resources.ReservationResourceSchema(ctx)
}

func (r *reservationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resources.ReservationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(reservationResourcePost(ctx, r.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *reservationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resources.ReservationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(reservationResourceGet(ctx, r.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *reservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resources.ReservationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// resp.Diagnostics.Append(reservationResourcePost(ctx, r.client, &data)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *reservationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resources.ReservationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(reservationResourceDelete(ctx, r.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
func (r *reservationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}
func (r *reservationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	// Retrieve import ID and save to id attribute
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("block"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[2])...)
}

func reservationResourceDelete(ctx context.Context, c *client.Client, data *resources.ReservationModel) diag.Diagnostics {
	var diags diag.Diagnostics
	requestData := client.ReservationApiModel{
		Space: data.Space.ValueString(),
		Block: data.Block.ValueString(),
		Id:    data.Id.ValueString(),
	}

	if err := c.ReservationApiDelete(ctx, requestData); err != nil {
		diags.AddError("Response Unmarshal Error", fmt.Sprintf("Failed to unmarshal response: %s", err))
		return diags
	}

	return diags
}

func reservationResourceGet(ctx context.Context, c *client.Client, data *resources.ReservationModel) diag.Diagnostics {
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

func reservationResourcePost(ctx context.Context, c *client.Client, data *resources.ReservationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	requestData := client.ReservationApiModel{
		Space:         data.Space.ValueString(),
		Block:         data.Block.ValueString(),
		SmallestCidr:  data.SmallestCidr.ValueBool(),
		ReverseSearch: data.ReverseSearch.ValueBool(),
		Size:          data.Size.ValueInt64(),
		Desc:          data.Desc.ValueString(),
		CIDR:          data.Cidr.ValueString(),
	}

	response, err := c.ReservationApiPost(ctx, requestData)
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
