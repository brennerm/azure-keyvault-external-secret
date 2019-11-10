package keyvault

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"k8s.io/klog"
)

var keyVaultClient keyvault.BaseClient
var url string

func getKeyvaultAuthorizer() (autorest.Authorizer, error) {
	// https://github.com/Azure/azure-sdk-for-go#more-authentication-details
	keyvaultAuthorizer, err := auth.NewAuthorizerFromEnvironment()
	return keyvaultAuthorizer, err
}

func Initialize(vaultUrl string) error {
	a, err := getKeyvaultAuthorizer()
	if err != nil {
		klog.Fatalf("Failed to get KeyVault Authorizer: %+v", err.Error())
	}

	keyVaultClient := keyvault.New()
	keyVaultClient.Authorizer = a

	url = vaultUrl

	return err
}

func GetSecret(name string, version string) string {
	ctx := context.Background()
	bundle, err := keyVaultClient.GetSecret(ctx, url, name, version)
	if err != nil {
		klog.Errorf("Failed to get secret %s: %+v", name, err.Error())
	}

	return *bundle.Value
}
