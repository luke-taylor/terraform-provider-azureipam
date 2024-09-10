package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Define the AdminsApiModel struct
type AdminsApiModel struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
}

// AdminsApiGet retrieves the list of admins and maps them to the AdminsModel.
func (c *Client) AdminsApiGet(ctx context.Context) ([]AdminsApiModel, error) {

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/admin/admins", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	// Execute the request
	respBody, err := c.DoRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response into a slice of AdminsApiModel
	var response []AdminsApiModel
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return response, nil
}
