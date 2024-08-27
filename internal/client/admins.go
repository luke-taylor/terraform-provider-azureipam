package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-azureipam/internal/gen/data_sources"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (c *Client) AdminsApiGet(ctx context.Context, admins *data_sources.AdminsModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/admin/admins", c.HostURL), nil)
	if err != nil {
		diags.AddError("Request Creation Error", fmt.Sprintf("Could not create HTTP request: %s", err))
		return diags
	}

	// Make the request
	adminsGet, err := c.DoRequest(req, &c.Token)
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
