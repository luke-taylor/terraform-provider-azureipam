package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *Client) ReservationsApiGet(ctx context.Context, requestData ReservationApiModel) ([]ReservationApiModel, error) {

	// Construct the URL for the GET request
	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations", c.HostURL, requestData.Space, requestData.Block)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Execute the request and obtain the response
	respBody, err := c.DoRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into a slice of ReservationApiModel
	var reservations []ReservationApiModel
	if err := json.Unmarshal(respBody, &reservations); err != nil {
		return nil, err
	}

	return reservations, nil
}

// ReservationApiGet handles GET requests for reservations
func (c *Client) ReservationApiGet(ctx context.Context, requestData ReservationApiModel) (*ReservationApiModel, error) {

	// Construct the URL for the GET request
	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations/%s", c.HostURL, requestData.Space, requestData.Block, requestData.Id)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Execute the request and obtain the response
	respBody, err := c.DoRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into a slice of ReservationApiModel
	var reservation = ReservationApiModel{}
	if err := json.Unmarshal(respBody, &reservation); err != nil {
		return nil, err
	}

	return &reservation, nil
}

// ReservationApiGetDelete handles GET and DELETE requests for reservations
func (c *Client) ReservationApiDelete(ctx context.Context, requestData ReservationApiModel) error {

	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations/%s", c.HostURL, requestData.Space, requestData.Block, requestData.Id)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	// Execute the request and obtain the response
	if _, err := c.DoRequest(req, &c.Token); err != nil {
		return err
	}

	return nil
}

// ReservationApiPost handles POST requests for reservations
func (c *Client) ReservationApiPost(ctx context.Context, requestData ReservationApiModel) (*ReservationApiModel, error) {

	url := fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations", c.HostURL, requestData.Space, requestData.Block)
	// Marshal the payload to JSON
	body, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	respBody, err := c.DoRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	reservation := ReservationApiModel{}
	if err := json.Unmarshal(respBody, &reservation); err != nil {
		return nil, err
	}

	return &reservation, nil
}
