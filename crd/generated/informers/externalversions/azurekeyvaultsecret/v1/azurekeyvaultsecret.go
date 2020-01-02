/*
 */

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	azurekeyvaultsecretv1 "github.com/brennerm/azure-keyvault-external-secret/crd/azurekeyvaultsecret/v1"
	versioned "github.com/brennerm/azure-keyvault-external-secret/crd/generated/clientset/versioned"
	internalinterfaces "github.com/brennerm/azure-keyvault-external-secret/crd/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/brennerm/azure-keyvault-external-secret/crd/generated/listers/azurekeyvaultsecret/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// AzureKeyVaultSecretInformer provides access to a shared informer and lister for
// AzureKeyVaultSecrets.
type AzureKeyVaultSecretInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.AzureKeyVaultSecretLister
}

type azureKeyVaultSecretInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewAzureKeyVaultSecretInformer constructs a new informer for AzureKeyVaultSecret type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewAzureKeyVaultSecretInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredAzureKeyVaultSecretInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredAzureKeyVaultSecretInformer constructs a new informer for AzureKeyVaultSecret type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredAzureKeyVaultSecretInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AzurekeyvaultsecretV1().AzureKeyVaultSecrets(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AzurekeyvaultsecretV1().AzureKeyVaultSecrets(namespace).Watch(options)
			},
		},
		&azurekeyvaultsecretv1.AzureKeyVaultSecret{},
		resyncPeriod,
		indexers,
	)
}

func (f *azureKeyVaultSecretInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredAzureKeyVaultSecretInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *azureKeyVaultSecretInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&azurekeyvaultsecretv1.AzureKeyVaultSecret{}, f.defaultInformer)
}

func (f *azureKeyVaultSecretInformer) Lister() v1.AzureKeyVaultSecretLister {
	return v1.NewAzureKeyVaultSecretLister(f.Informer().GetIndexer())
}