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

type reservationApiModel struct {
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

func (c *Client) ReservationsApiGet(ctx context.Context, data *data_sources.ReservationsModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Construct the URL for the GET request
	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations",
		c.HostURL, strings.Trim(data.Space.ValueString(), "\""), strings.Trim(data.Block.ValueString(), "\""))

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		diags.AddError("Request Creation Error", fmt.Sprintf("Could not create HTTP request: %s", err))
		return diags
	}

	// Execute the request and obtain the response
	respBody, err := c.DoRequest(req, &c.Token)
	if err != nil {
		diags.AddError("API Request Error", fmt.Sprintf("API request failed: %s", err))
		return diags
	}

	// Unmarshal the JSON response into a slice of reservationApiModel
	var reservations []reservationApiModel
	if err := json.Unmarshal(respBody, &reservations); err != nil {
		diags.AddError("Response Unmarshal Error", fmt.Sprintf("Failed to unmarshal response: %s", err))
		return diags
	}
	// print the response
	elements := make([]attr.Value, len(reservations))
	for i, reservation := range reservations {
		tags := make(map[string]attr.Value, len(reservation.Tag))
		for k, v := range reservation.Tag {
			tags[k] = types.StringValue(v)
		}
		createdOn := big.NewFloat(reservation.CreatedOn)
		settledOn := big.NewFloat(reservation.SettledOn)
		tag, _ := types.MapValue(types.StringType, tags)
		objVal, objDiags := data_sources.NewReservationsValue(data_sources.NewReservationsValueNull().AttributeTypes(ctx),
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
		diags.Append(objDiags...)
		if diags.HasError() {
			return diags
		}
		elements[i] = objVal
	}
	fmt.Println("here")
	fmt.Println(elements)

	// Set the Reservations field in the ReservationsModel
	data.Reservations, diags = types.SetValueFrom(ctx, data_sources.NewReservationsValueNull().Type(ctx), &elements)

	return diags
}

// ReservationApiGet handles GET requests for reservations
func (c *Client) ReservationApiGet(ctx context.Context, data *data_sources.ReservationModel) diag.Diagnostics {
	payload := reservationApiModel{
		Space: data.Space.ValueString(),
		Block: data.Block.ValueString(),
		Id:    data.Id.ValueString(),
	}

	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations/%s",
		c.HostURL, strings.Trim(data.Space.ValueString(), "\""), strings.Trim(data.Block.ValueString(), "\""), strings.Trim(data.Id.ValueString(), "\""))

	response, diags := c.reservationExecuteRequest(ctx, "GET", url, payload)
	if diags.HasError() {
		return diags
	}

	return mapApiResponseToModel(response, data)
}

// ReservationApiGetDelete handles GET and DELETE requests for reservations
func (c *Client) ReservationApiGetDelete(ctx context.Context, data *resources.ReservationModel, method string) diag.Diagnostics {
	payload := reservationApiModel{
		Space: data.Space.ValueString(),
		Block: data.Block.ValueString(),
		Id:    data.Id.ValueString(),
	}

	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations/%s",
		c.HostURL, strings.Trim(data.Space.ValueString(), "\""), strings.Trim(data.Block.ValueString(), "\""), strings.Trim(data.Id.ValueString(), "\""))

	if method == "DELETE" {
		_, diags := c.reservationExecuteRequest(ctx, method, url, payload)
		return diags
	}

	response, diags := c.reservationExecuteRequest(ctx, method, url, payload)
	if diags.HasError() {
		return diags
	}

	return mapApiResponseToModel(response, data)
}

// ReservationApiPost handles POST requests for reservations
func (c *Client) ReservationApiPost(ctx context.Context, data *resources.ReservationModel) diag.Diagnostics {
	payload := reservationApiModel{
		Space:         data.Space.ValueString(),
		Block:         data.Block.ValueString(),
		SmallestCidr:  data.SmallestCidr.ValueBool(),
		ReverseSearch: data.ReverseSearch.ValueBool(),
		Size:          data.Size.ValueInt64(),
		Desc:          data.Desc.ValueString(),
		CIDR:          data.Cidr.ValueString(),
	}

	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations", c.HostURL, strings.Trim(data.Space.ValueString(), "\""), strings.Trim(data.Block.ValueString(), "\""))

	response, diags := c.reservationExecuteRequest(ctx, "POST", url, payload)
	if diags.HasError() {
		return diags
	}

	return mapApiResponseToModel(response, data)
}

// reservationExecuteRequest handles making the HTTP request and unmarshalling the response
func (c *Client) reservationExecuteRequest(ctx context.Context, method, url string, payload reservationApiModel) (reservationApiModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Marshal the payload to JSON
	reservationData, err := json.Marshal(payload)
	if err != nil {
		diags.AddError("Failed to marshal reservation data", err.Error())
		return reservationApiModel{}, diags
	}

	// Create the HTTP request with context
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reservationData))
	if err != nil {
		diags.AddError("Failed to create HTTP request", err.Error())
		return reservationApiModel{}, diags
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	respBody, err := c.DoRequest(req, &c.Token)
	if err != nil {
		select {
		case <-ctx.Done(): // Handle context cancellation or timeout
			diags.AddError("Request canceled or timed out", ctx.Err().Error())
		default:
			diags.AddError("API request failed", err.Error())
		}
		return reservationApiModel{}, diags
	}

	// Initialize the response model
	response := reservationApiModel{}

	// Only unmarshal the response if the method is not DELETE
	if method != "DELETE" {
		if err := json.Unmarshal(respBody, &response); err != nil {
			diags.AddError("Failed to unmarshal API response", err.Error())
			return reservationApiModel{}, diags
		}
	}

	return response, diags
}

// Helper function to map API response to Terraform model
func mapApiResponseToModel(response reservationApiModel, model interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	switch data := model.(type) {
	case *data_sources.ReservationModel:
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

	case *resources.ReservationModel:
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

		if response.Size != 0 {
			data.Size = types.Int64Value(response.Size)
		}

		if response.Tag != nil {
			tagElements := make(map[string]attr.Value)
			for k, v := range response.Tag {
				tagElements[k] = types.StringValue(v)
			}
			data.Tag, _ = types.MapValue(types.StringType, tagElements)
		} else {
			data.Tag = types.MapNull(types.StringType)
		}
	}

	return diags
}
