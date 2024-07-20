package provider

import (
    "context"
    "terraform-provider-petstore/internal/datasource_admins"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = (*adminsDataSource)(nil)

func NewAdminsDataSource() datasource.DataSource {
    return &adminsDataSource{}
}

type adminsDataSource struct{}

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

    resp.Diagnostics.Append(callAdminsAPI(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Typically this method would contain logic that makes an HTTP call to a remote API, and then stores
// computed results back to the data model. For example purposes, this function just sets computed Admins
// values to mock values to avoid data consistency errors.
func callAdminsAPI(ctx context.Context, admins *datasource_admins.AdminsModel) diag.Diagnostics {
    admins.PetId = types.Int64Value(1)
    admins.Quantity = types.Int64Value(10)
    admins.Status = types.StringValue("delivered")
    admins.ShipDate = types.StringValue("2023-07-25T23:43:16Z")
    admins.Complete = types.BoolValue(true)

    return nil
}
