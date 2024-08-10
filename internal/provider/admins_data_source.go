package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-azureipam/internal/client"
	"terraform-provider-azureipam/internal/datasource_admins"

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

type AdminsValueRaw struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
}


func (d *adminsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_admins"
}

func (d *adminsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_admins.AdminsDataSourceSchema(ctx)
}

func (d *adminsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_admins.AdminsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(callAdminsAPI(ctx, &data, d.client)...)
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
func callAdminsAPI(ctx context.Context, admins *datasource_admins.AdminsModel, client *client.Client) diag.Diagnostics {
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

	// Unmarshal the JSON response into the adminsData slice using native Go types
	var adminsDataRaw []AdminsValueRaw
	err = json.Unmarshal(adminsGet, &adminsDataRaw)
	if err != nil {
			diags.AddError("Response Unmarshal Error", fmt.Sprintf("Failed to unmarshal response: %s", err))
			return diags
	}

	// Convert the raw data to Terraform types
	adminsData := make([]datasource_admins.AdminsValue, len(adminsDataRaw))
	for i, adminRaw := range adminsDataRaw {
			adminsData[i] = datasource_admins.AdminsValue{
					Id:         types.StringValue(adminRaw.Id),
					Name:       types.StringValue(adminRaw.Name),
					Email:      types.StringValue(adminRaw.Email),
					AdminsType: types.StringValue(adminRaw.Type),
			}
	}

	// Prepare elements for the Set
	elements := make([]attr.Value, len(adminsData))
	for i, admin := range adminsData {
			objVal, objDiags := types.ObjectValue(
					map[string]attr.Type{
							"id":    types.StringType,
							"name":  types.StringType,
							"email": types.StringType,
							"type":  types.StringType,
					},
					map[string]attr.Value{
							"id":    admin.Id,
							"name":  admin.Name,
							"email": admin.Email,
							"type":  admin.AdminsType,
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

