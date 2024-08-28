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

// Define the AdminsApiModel struct
type adminsApiModel struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
}

// AdminsApiGet retrieves the list of admins and maps them to the AdminsModel.
func (c *Client) AdminsApiGet(ctx context.Context, data *data_sources.AdminsModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/admin/admins", c.HostURL), nil)
	if err != nil {
		diags.AddError("Request Creation Error", fmt.Sprintf("Could not create HTTP request: %s", err))
		return diags
	}

	// Execute the request
	respBody, err := c.DoRequest(req, &c.Token)
	if err != nil {
		diags.AddError("API Request Error", fmt.Sprintf("API request failed: %s", err))
		return diags
	}

	// Unmarshal the JSON response into a slice of adminsApiModel
	var response []adminsApiModel
	if err := json.Unmarshal(respBody, &response); err != nil {
		diags.AddError("Response Unmarshal Error", fmt.Sprintf("Failed to unmarshal response: %s", err))
		return diags
	}

	// Convert the response to Terraform types
	elements := make([]attr.Value, len(response))
	for i, admin := range response {
		objVal, objDiags := types.ObjectValue(
			map[string]attr.Type{
				"id":    types.StringType,
				"name":  types.StringType,
				"email": types.StringType,
				"type":  types.StringType,
			},
			map[string]attr.Value{
				"id":    types.StringValue(admin.ID),
				"name":  types.StringValue(admin.Name),
				"email": types.StringValue(admin.Email),
				"type":  types.StringValue(admin.Type),
			},
		)
		diags.Append(objDiags...)
		if diags.HasError() {
			return diags
		}
		elements[i] = objVal
	}

	// Set the Admins field in the AdminsModel
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

	data.Admins = adminsSet

	return diags
}
