package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"terraform-provider-azureipam/internal/gen/data_sources"
	"terraform-provider-azureipam/internal/gen/resources"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ReservationApiModel struct {
	Id            string            `json:"id,omitempty"`
	Space         string            `json:"space,omitempty"`
	Block         string            `json:"block,omitempty"`
	CIDR          string            `json:"cidr,omitempty"`
	Desc          string            `json:"desc,omitempty"`
	CreatedOn     float64           `json:"createdOn,omitempty"`
	CreatedBy     string            `json:"createdBy,omitempty"`
	SettledBy     string            `json:"settledBy,omitempty"`
	SettledOn     float64           `json:"settledOn,omitempty"`
	Status        string            `json:"status,omitempty"`
	Tag           map[string]string `json:"tag,omitempty"`
	ReverseSearch bool              `json:"reverse_search,omitempty"`
	Size          int64             `json:"size,omitempty"`
	SmallestCidr  bool              `json:"smallest_cidr,omitempty"`
}

func (c *Client) ReservationApiGet(ctx context.Context, data *data_sources.ReservationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	payload := ReservationApiModel{
		Space: data.Space.ValueString(),
		Block: data.Block.ValueString(),
		Id:    data.Id.ValueString(),
	}

	// Construct the API URL
	space := strings.Trim(data.Space.ValueString(), "\"")
	block := strings.Trim(data.Block.ValueString(), "\"")
	id := strings.Trim(data.Id.ValueString(), "\"")
	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations/%s", c.HostURL, space, block, id)

	// Marshal the payload to JSON
	reservationData, err := json.Marshal(payload)
	if err != nil {
		diags.AddError("Failed to marshal reservation data", err.Error())
		return diags
	}

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(reservationData))
	if err != nil {
		diags.AddError("Failed to create HTTP request", err.Error())
		return diags
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	respBody, err := c.DoRequest(req, &c.Token)
	if err != nil {
		diags.AddError("API request failed", err.Error())
		return diags
	}
	// Assuming `DoRequest` returns raw JSON data and status code together
	// Unmarshal the response into a temporary struct that matches the API response
	response := ReservationApiModel{}
	// Unmarshal the response body
	if err := json.Unmarshal(respBody, &response); err != nil {
		diags.AddError("Failed to unmarshal API response", err.Error())
		return diags
	}
	// add write respones to standard output + add HERE to make more idetnifable

	// Map the API response to the Terraform model
	settledOnBigFloat := big.NewFloat(response.SettledOn)
	data.SettledOn = types.NumberValue(settledOnBigFloat)
	data.SettledBy = types.StringValue(response.SettledBy)
	data.Status = types.StringValue(response.Status)
	data.Id = types.StringValue(response.Id)
	data.Cidr = types.StringValue(response.CIDR)
	data.CreatedBy = types.StringValue(response.CreatedBy)
	data.Desc = types.StringValue(response.Desc)
	createdOnBigFloat := big.NewFloat(response.CreatedOn)
	data.CreatedOn = types.NumberValue(createdOnBigFloat)
	data.Status = types.StringValue(response.Status)
	if response.Tag != nil {
		tagElements := make(map[string]attr.Value)
		for k, v := range response.Tag {
			tagElements[k] = types.StringValue(v)
		}
		data.Tag, _ = types.MapValue(types.StringType, tagElements)
		if err != nil {
			diags.AddError("Failed to create tag map", err.Error())
		}
	} else {
		data.Tag = types.MapNull(types.StringType)
	}

	return diags
}

func (c *Client) ReservationApiGetDelete(ctx context.Context, data *resources.ReservationModel, method string) diag.Diagnostics {
	var diags diag.Diagnostics

	payload := ReservationApiModel{
		Space: data.Space.ValueString(),
		Block: data.Block.ValueString(),
		Id:    data.Id.ValueString(),
	}

	// Construct the API URL
	space := strings.Trim(data.Space.ValueString(), "\"")
	block := strings.Trim(data.Block.ValueString(), "\"")
	id := strings.Trim(data.Id.ValueString(), "\"")
	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations/%s", c.HostURL, space, block, id)

	// Marshal the payload to JSON
	reservationData, err := json.Marshal(payload)
	if err != nil {
		diags.AddError("Failed to marshal reservation data", err.Error())
		return diags
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reservationData))
	if err != nil {
		diags.AddError("Failed to create HTTP request", err.Error())
		return diags
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	respBody, err := c.DoRequest(req, &c.Token)
	if err != nil {
		diags.AddError("API request failed", err.Error())
		return diags
	}
	// Assuming `DoRequest` returns raw JSON data and status code together
	// Unmarshal the response into a temporary struct that matches the API response
	if method != "DELETE" {
		response := ReservationApiModel{}
		// Unmarshal the response body
		if err := json.Unmarshal(respBody, &response); err != nil {
			diags.AddError("Failed to unmarshal API response", err.Error())
			return diags
		}
		// add write respones to standard output + add HERE to make more idetnifable

		// Map the API response to the Terraform model
		settledOnBigFloat := big.NewFloat(response.SettledOn)
		data.SettledOn = types.NumberValue(settledOnBigFloat)
		data.SettledBy = types.StringValue(response.SettledBy)
		data.Status = types.StringValue(response.Status)
		data.Id = types.StringValue(response.Id)
		data.Cidr = types.StringValue(response.CIDR)
		data.CreatedBy = types.StringValue(response.CreatedBy)
		data.Desc = types.StringValue(response.Desc)
		createdOnBigFloat := big.NewFloat(response.CreatedOn)
		data.CreatedOn = types.NumberValue(createdOnBigFloat)
		data.Status = types.StringValue(response.Status)
		if data.Size.ValueInt64() == 0 {
			data.Size = types.Int64Value(response.Size)
		} // data.Size = types.Int64Value(response.Size)
		if response.Tag != nil {
			tagElements := make(map[string]attr.Value)
			for k, v := range response.Tag {
				tagElements[k] = types.StringValue(v)
			}
			data.Tag, _ = types.MapValue(types.StringType, tagElements)
			if err != nil {
				diags.AddError("Failed to create tag map", err.Error())
			}
		} else {
			data.Tag = types.MapNull(types.StringType)
		}
	}
	return diags
}

func (c *Client) ReservationApiPost(ctx context.Context, data *resources.ReservationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	payload := ReservationApiModel{
		Space:         data.Space.ValueString(),
		Block:         data.Block.ValueString(),
		SmallestCidr:  data.SmallestCidr.ValueBool(),
		ReverseSearch: data.ReverseSearch.ValueBool(),
		Size:          data.Size.ValueInt64(),
		Desc:          data.Desc.ValueString(),
		CIDR:          data.Cidr.ValueString(),
	}

	// Construct the API URL
	space := strings.Trim(data.Space.ValueString(), "\"")
	block := strings.Trim(data.Block.ValueString(), "\"")
	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations", c.HostURL, space, block)

	// Marshal the payload to JSON
	reservationData, err := json.Marshal(payload)
	if err != nil {
		diags.AddError("Failed to marshal reservation data", err.Error())
		return diags
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reservationData))
	if err != nil {
		diags.AddError("Failed to create HTTP request", err.Error())
		return diags
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	respBody, err := c.DoRequest(req, &c.Token)
	if err != nil {
		diags.AddError("API request failed", err.Error())
		return diags
	}

	// Assuming `DoRequest` returns raw JSON data and status code together
	// Unmarshal the response into a temporary struct that matches the API response

	response := ReservationApiModel{}
	// Unmarshal the response body
	if err := json.Unmarshal(respBody, &response); err != nil {
		diags.AddError("Failed to unmarshal API response", err.Error())
		return diags
	}
	// Map the API response to the Terraform model
	data.Id = types.StringValue(response.Id)
	data.Cidr = types.StringValue(response.CIDR)
	data.CreatedBy = types.StringValue(response.CreatedBy)
	data.Desc = types.StringValue(response.Desc)
	settledOnBigFloat := big.NewFloat(response.SettledOn)
	data.SettledOn = types.NumberValue(settledOnBigFloat)
	data.SettledBy = types.StringValue(response.SettledBy)
	createdOnBigFloat := big.NewFloat(response.CreatedOn)
	data.CreatedOn = types.NumberValue(createdOnBigFloat)
	data.Status = types.StringValue(response.Status)
	if data.Size.ValueInt64() == 0 {
		data.Size = types.Int64Value(response.Size)
	} // data.Size = types.Int64Value(response.Size)
	if response.Tag != nil {
		tagElements := make(map[string]attr.Value)
		for k, v := range response.Tag {
			tagElements[k] = types.StringValue(v)
		}
		data.Tag, _ = types.MapValue(types.StringType, tagElements)
		if err != nil {
			diags.AddError("Failed to create tag map", err.Error())
		}
	} else {
		data.Tag = types.MapNull(types.StringType)
	}
	return diags
}
