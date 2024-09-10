package provider

import (
	"context"
	"fmt"
	"terraform-provider-azureipam/internal/client"
	"terraform-provider-azureipam/internal/gen/data_sources"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*adminsDataSource)(nil)

func NewAdminsDataSource() datasource.DataSource {
	return &adminsDataSource{}
}

type adminsDataSource struct {
	client *client.Client
}

func (d *adminsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_admins"
}

func (d *adminsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = data_sources.AdminsDataSourceSchema(ctx)
}

func (d *adminsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data data_sources.AdminsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(adminsDataGet(ctx, d.client, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
func (d *adminsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func adminsDataGet(ctx context.Context, c *client.Client, data *data_sources.AdminsModel) diag.Diagnostics {
	var diags diag.Diagnostics

	response, err := c.AdminsApiGet(ctx)
	if err != nil {
		diags.AddError("Response Unmarshal Error", fmt.Sprintf("Failed to unmarshal response: %s", err))
		return diags
	}

	elements := make([]attr.Value, len(response))
	for i, admin := range response {
		elementValue := data_sources.NewAdminsValueMust(
			data_sources.NewAdminsValueNull().AttributeTypes(ctx),
			map[string]attr.Value{
				"id":    types.StringValue(admin.ID),
				"email": types.StringValue(admin.Email),
				"name":  types.StringValue(admin.Name),
				"type":  types.StringValue(admin.Type),
			},
		)
		elements[i] = elementValue
	}

	data.Admins, diags = types.SetValueFrom(ctx, data_sources.NewAdminsValueNull().Type(ctx), &elements)

	return diags

}
