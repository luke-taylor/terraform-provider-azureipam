package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"terraform-provider-azureipam/internal/client"
	"terraform-provider-azureipam/internal/resource_reservation"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ReservationAPIModel struct {
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
	Reservation   string            `json:"reservation,omitempty"`
	SmallestCidr  bool              `json:"smallest_cidr,omitempty"`
}

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
	resp.Schema = resource_reservation.ReservationResourceSchema(ctx)
}

func (r *reservationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_reservation.ReservationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(callReservationAPI(ctx, &data, r.client, "POST")...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *reservationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_reservation.ReservationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(callReservationAPIGet(ctx, &data, r.client, "GET")...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *reservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_reservation.ReservationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(callReservationAPI(ctx, &data, r.client, "POST")...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *reservationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_reservation.ReservationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(callReservationAPIGet(ctx, &data, r.client, "DELETE")...)
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
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func callReservationAPIGet(ctx context.Context, data *resource_reservation.ReservationModel, client *client.Client, method string) diag.Diagnostics {
	var diags diag.Diagnostics

	payload := ReservationAPIModel{
		Space: data.Space.ValueString(),
		Block: data.Block.ValueString(),
		Id:    data.Id.ValueString(),
	}

	// Construct the API URL
	space := strings.Trim(data.Space.ValueString(), "\"")
	block := strings.Trim(data.Block.ValueString(), "\"")
	id := strings.Trim(data.Id.ValueString(), "\"")
	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations/%s", client.HostURL, space, block, id)

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
	respBody, err := client.DoRequest(req, &client.Token)
	if err != nil {
		diags.AddError("API request failed", err.Error())
		return diags
	}

	// Assuming `DoRequest` returns raw JSON data and status code together
	// Unmarshal the response into a temporary struct that matches the API response
	if method != "DELETE" {
		response := ReservationAPIModel{}
		// Unmarshal the response body
		if err := json.Unmarshal(respBody, &response); err != nil {
			diags.AddError("Failed to unmarshal API response", err.Error())
			return diags
		}

		// Map the API response to the Terraform model
		data.Id = types.StringValue(response.Id)
		data.Cidr = types.StringValue(response.CIDR)
		data.CreatedBy = types.StringValue(response.CreatedBy)
		data.ReverseSearch = types.BoolValue(response.ReverseSearch)
		data.Size = types.Int64Value(response.Size)
		data.Desc = types.StringValue(response.Desc)
		settledOnBigFloat := big.NewFloat(response.SettledOn)
		data.SettledOn = types.NumberValue(settledOnBigFloat)
		data.SettledBy = types.StringValue(response.SettledBy)
		createdOnBigFloat := big.NewFloat(response.CreatedOn)
		data.CreatedOn = types.NumberValue(createdOnBigFloat)
		data.Status = types.StringValue(response.Status)
		data.Reservation = types.StringValue(response.Reservation)
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

func callReservationAPI(ctx context.Context, data *resource_reservation.ReservationModel, client *client.Client, method string) diag.Diagnostics {
	var diags diag.Diagnostics

	payload := ReservationAPIModel{
		Space:         data.Space.ValueString(),
		Block:         data.Block.ValueString(),
		SmallestCidr:  data.SmallestCidr.ValueBool(),
		ReverseSearch: data.ReverseSearch.ValueBool(),
		Size:          data.Size.ValueInt64(),
	}

	// Construct the API URL
	space := strings.Trim(data.Space.ValueString(), "\"")
	block := strings.Trim(data.Block.ValueString(), "\"")
	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations", client.HostURL, space, block)

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
	respBody, err := client.DoRequest(req, &client.Token)
	if err != nil {
		diags.AddError("API request failed", err.Error())
		return diags
	}

	// Assuming `DoRequest` returns raw JSON data and status code together
	// Unmarshal the response into a temporary struct that matches the API response

	response := ReservationAPIModel{}
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
	data.Reservation = types.StringValue(response.Reservation)
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
