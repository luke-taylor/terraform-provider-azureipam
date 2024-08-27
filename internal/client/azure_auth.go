package client

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func GetAzureAccessToken(apiGuid string) (string, error) {
	// Create a credential using the default Azure credential chain
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return "", fmt.Errorf("failed to obtain a credential: %v", err)
	}

	// Create a context
	ctx := context.Background()

	// Request a token
	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{fmt.Sprintf("api://%s/.default", apiGuid)},
	})
	if err != nil {
		return "", fmt.Errorf("failed to get token: %v", err)
	}

	bearerToken := fmt.Sprintf("Bearer %s", token.Token)
	return bearerToken, nil
}
