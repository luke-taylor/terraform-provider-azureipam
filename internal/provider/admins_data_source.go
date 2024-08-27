package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-azureipam/internal/client"
	"terraform-provider-azureipam/internal/gen"

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
	resp.Schema = gen.AdminsDataSourceSchema(ctx)
}

func (d *adminsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data gen.AdminsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(adminsApiGet(ctx, &data, d.client)...)
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
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Typically this method would contain logic that makes an HTTP call to a remote API, and then stores
// computed results back to the data model. For example purposes, this function just sets computed Admins
// values to mock values to avoid data consistency errors.
func adminsApiGet(ctx context.Context, admins *gen.AdminsModel, client *client.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/admin/admins", client.HostURL), nil)
	if err != nil {
		diags.AddError("Request Creation Error", fmt.Sprintf("Could not create HTTP request: %s", err))
		return diags
	}

	// Make the request
	adminsGet, err := client.DoRequest(req, &client.Token)
	if err != nil {
		diags.AddError("API Request Error", fmt.Sprintf("API request failed: %s", err))
		return diags
	}

	// Unmarshal the JSON response into a slice of maps
	var adminsDataRaw []map[string]interface{}
	err = json.Unmarshal(adminsGet, &adminsDataRaw)
	if err != nil {
		diags.AddError("Response Unmarshal Error", fmt.Sprintf("Failed to unmarshal response: %s", err))
		return diags
	}

	// Convert the raw data to Terraform types
	elements := make([]attr.Value, len(adminsDataRaw))
	for i, adminRaw := range adminsDataRaw {
		// Convert each map entry to Terraform types
		objVal, objDiags := types.ObjectValue(
			map[string]attr.Type{
				"id":    types.StringType,
				"name":  types.StringType,
				"email": types.StringType,
				"type":  types.StringType,
			},
			map[string]attr.Value{
				"id":    types.StringValue(adminRaw["id"].(string)),
				"name":  types.StringValue(adminRaw["name"].(string)),
				"email": types.StringValue(adminRaw["email"].(string)),
				"type":  types.StringValue(adminRaw["type"].(string)),
			},
		)
		diags.Append(objDiags...)
		if diags.HasError() {
			return diags
		}
		elements[i] = objVal
	}

	// Set the admins value
	adminsSet, setDiags := types.SetValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":    types.StringType,
				"name":  types.StringType,
				"email": types.StringType,
				"type":  types.StringType,
			},
		},
		elements,
	)
	diags.Append(setDiags...)
	if diags.HasError() {
		return diags
	}

	admins.Admins = adminsSet

	return diags
}
