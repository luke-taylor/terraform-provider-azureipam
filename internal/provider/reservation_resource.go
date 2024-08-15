package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"terraform-provider-azureipam/internal/client"
	"terraform-provider-azureipam/internal/resource_reservation"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

	// Construct the API request URL

	// Create a new struct that includes Space and Block
	type ReservationPayload struct {
		Space        string `json:"space"`
		Block        string `json:"block"`
		SmallestCidr bool   `json:"smallest_cidr"`
		Size         string `json:"size"`
	}

	payload := ReservationPayload{
		Space:        data.Space.String(),
		Block:        data.Block.String(),
		SmallestCidr: data.SmallestCidr.ValueBool(),
		Size:         "24",
	}
	space := strings.Trim(data.Space.String(), "\"")
	block := strings.Trim(data.Block.String(), "\"")

	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations", client.HostURL, space, block)
	// Marshal the reservation data to JSON
	reservationData, err := json.Marshal(payload)
	if err != nil {
		diags.AddError("Failed to marshal reservation data", err.Error())
		return diags
	}

	// Debugging: Log the JSON payload
	fmt.Printf("HERE...Request Body: %s\n", string(reservationData)+url)

	// Create a new HTTP request with the specified method
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reservationData))
	if err != nil {
		diags.AddError("Failed to create HTTP request", err.Error())
		return diags
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request using the client's DoRequest method
	respBody, err := client.DoRequest(req, &client.Token)
	if err != nil {
		diags.AddError("API request failed", err.Error())
		return diags
	}

	// Log the response body for debugging
	fmt.Printf("Response Body: %s\n", string(respBody))

	// Parse the response and update the model
	if err := json.Unmarshal(respBody, data); err != nil {
		diags.AddError("Failed to unmarshal API response", err.Error())
		return diags
	}

	return diags
}
