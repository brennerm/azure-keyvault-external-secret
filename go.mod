module github.com/brennerm/azure-keyvault-external-secret

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v37.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.2
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.1
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	k8s.io/api v0.0.0-20191121015604-11707872ac1c
	k8s.io/apiextensions-apiserver v0.0.0-20191204090421-cd61debedab5
	k8s.io/apimachinery v0.0.0-20191203211716-adc6f4cd9e7d
	k8s.io/client-go v0.0.0-20191204082520-bc9b51d240b2
	k8s.io/code-generator v0.0.0-20191121015212-c4c8f8345c7e
	k8s.io/klog v1.0.0
	k8s.io/sample-controller v0.0.0-20191123021055-65e4a1173959
)
