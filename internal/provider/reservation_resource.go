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

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

// Typically this method would contain logic that makes an HTTP call to a remote API, and then stores
// computed results back to the data model. For example purposes, this function just sets all unknown
// Reservation values to null to avoid data consistency errors.
func callReservationAPI(ctx context.Context, data *resource_reservation.ReservationModel, client *client.Client, method string) diag.Diagnostics {
	var diags diag.Diagnostics

	// Define the payload based on the input data
	type ReservationPayload struct {
		Space        string `json:"space"`
		Block        string `json:"block"`
		SmallestCidr bool   `json:"smallest_cidr"`
		Size         string `json:"size"`
	}

	payload := ReservationPayload{
		Space:        data.Space.ValueString(),
		Block:        data.Block.ValueString(),
		SmallestCidr: data.SmallestCidr.ValueBool(),
		Size:         "24",
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
	var response struct {
		ID            string            `json:"id"`
		Space         string            `json:"space"`
		Block         string            `json:"block"`
		CIDR          string            `json:"cidr"`
		Desc          *string           `json:"desc"`
		CreatedOn     float64           `json:"createdOn"`
		CreatedBy     string            `json:"createdBy"`
		SettledBy     *string           `json:"settledBy"`
		SettledOn     *float64          `json:"settledOn"`
		Status        string            `json:"status"`
		Tag           map[string]string `json:"tag"`
		XIPAMResID    string            `json:"x_ipam_res_id,omitempty"`
		ReverseSearch bool              `json:"reverse_search"`
		Size          int64             `json:"size"`
		Reservation   string            `json:"reservation"`
	}

	// Unmarshal the response body
	if err := json.Unmarshal(respBody, &response); err != nil {
		diags.AddError("Failed to unmarshal API response", err.Error())
		return diags
	}

	// Map the API response to the Terraform model
	data.Id = types.StringValue(response.ID)
	data.Space = types.StringValue(response.Space)
	data.Block = types.StringValue(response.Block)
	data.Cidr = types.StringValue(response.CIDR)
	data.CreatedBy = types.StringValue(response.CreatedBy)
	data.ReverseSearch = types.BoolValue(response.ReverseSearch)
	data.Size = types.Int64Value(response.Size)
	data.Reservation = types.StringValue(response.Reservation)

	// Convert float64 to *big.Float for CreatedOn
	createdOnBigFloat := big.NewFloat(response.CreatedOn)
	data.CreatedOn = types.NumberValue(createdOnBigFloat)

	data.Status = types.StringValue(response.Status)

	// Handle optional fields
	if response.Desc != nil {
		data.Desc = types.StringValue(*response.Desc)
	} else {
		data.Desc = types.StringNull()
	}

	if response.SettledBy != nil {
		data.SettledBy = types.StringValue(*response.SettledBy)
	} else {
		data.SettledBy = types.StringNull()
	}

	if response.SettledOn != nil {
		settledOnBigFloat := big.NewFloat(*response.SettledOn)
		data.SettledOn = types.NumberValue(settledOnBigFloat)
	} else {
		data.SettledOn = types.NumberNull()
	}

	// Handling the Tag field
	if response.Tag != nil {
		// Prepare the elements slice
		// elements := make(map[string]attr.Value, len(response.Tag))

		// for k, v := range response.Tag {
		// 	elements[k] = types.StringValue(v)
		// 	// Check if the Tag key is "X-IPAM-RES-ID" and set the XIPAMResID in data
		// 	if k == "X-IPAM-RES-ID" {

		// 	}
		// }

		// // Convert the map to an ObjectValue and set it in the data model
		data.Tag = resource_reservation.NewTagValueNull()

	} else {
		data.Tag = resource_reservation.NewTagValueNull()
	}
	return diags
}
