package keyvault

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"k8s.io/klog"
)

type KeyVaultClient struct {
	client keyvault.BaseClient
}

func getKeyvaultAuthorizer() (autorest.Authorizer, error) {
	// https://github.com/Azure/azure-sdk-for-go#more-authentication-details
	keyvaultAuthorizer, err := auth.NewAuthorizerFromEnvironmentWithResource("https://vault.azure.net")
	return keyvaultAuthorizer, err
}

func NewKeyVaultClient() (KeyVaultClient, error) {
	a, err := getKeyvaultAuthorizer()
	if err != nil {
		klog.Fatalf("Failed to get KeyVault Authorizer: %+v", err.Error())
	}

	keyVaultClient := KeyVaultClient{
		client: keyvault.New(),
	}

	keyVaultClient.client.Authorizer = a

	return keyVaultClient, err
}

func (c *KeyVaultClient) GetSecret(vaultURL string, name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bundle, err := c.client.GetSecret(ctx, vaultURL, name, "") // we'll always get the latest version
	if err != nil {
		return "", fmt.Errorf("Failed to get secret %s: %+v", name, err.Error())
	}

	return *bundle.Value, nil
}

func (c *KeyVaultClient) GetSecrets(vaultURL string, names []string) ([]string, error) {
	var result []string

	for _, name := range names {
		value, err := c.GetSecret(vaultURL, name)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	return result, nil
}
