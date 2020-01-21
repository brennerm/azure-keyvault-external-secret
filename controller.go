package main

import (
	"fmt"
	"time"

	clientset "github.com/brennerm/azure-keyvault-external-secret/crd/generated/clientset/versioned"
	informers "github.com/brennerm/azure-keyvault-external-secret/crd/generated/informers/externalversions/azurekeyvaultsecret/v1"
	listers "github.com/brennerm/azure-keyvault-external-secret/crd/generated/listers/azurekeyvaultsecret/v1"
	"github.com/brennerm/azure-keyvault-external-secret/keyvault"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

// Controller is the controller implementation for Foo resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	sampleclientset clientset.Interface

	keyVaultClient keyvault.KeyVaultClient

	// K8s secrets
	secretsLister corelisters.SecretLister
	secretsSynced cache.InformerSynced
	// Azure Key Vault secrets
	kvSecretsLister listers.AzureKeyVaultSecretLister
	kvSecretsSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new Azure Keyvault secret controller
func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	secretInformer coreinformers.SecretInformer,
	kvSecretInformer informers.AzureKeyVaultSecretInformer,
	keyVaultClient keyvault.KeyVaultClient) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(scheme.AddToScheme(scheme.Scheme))
	klog.Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "azure-keyvault-secret-controller"})

	controller := &Controller{
		kubeclientset:   kubeclientset,
		sampleclientset: sampleclientset,
		keyVaultClient:  keyVaultClient,
		secretsLister:   secretInformer.Lister(),
		secretsSynced:   secretInformer.Informer().HasSynced,
		kvSecretsLister: kvSecretInformer.Lister(),
		kvSecretsSynced: kvSecretInformer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "AzureKeyVaultSecrets"),
		recorder:        recorder,
	}

	klog.Info("Setting up event handlers for AzureKeyVaultSecret")

	// Set up an event handler for when AzureKeyVaultSecret resources change
	kvSecretInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.pushKvSecret,
		UpdateFunc: func(old, new interface{}) {
			controller.pushKvSecret(new)
		},
		DeleteFunc: func(obj interface{}) {
			controller.pushKvSecret(obj)
		},
	})

	return controller
}

func (c *Controller) pushKvSecret(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		klog.Errorf("failed to generate cache key for object %+v: %+x", obj, err)
		return
	}

	c.workqueue.Add(key)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting azure-keyvault-secret controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.kvSecretsSynced, c.secretsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch workers to process resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string // key should contain $namespace/$name
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj) // remove key from queue to prevent processing it again
			klog.Errorf("Expected string in workqueue but got %T, removing from queue", obj)
			return nil
		}

		if err := c.handleKvSecretChange(key); err != nil {
			return fmt.Errorf("Error syncing '%s': %s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced \"%s\"", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// handles a creation, update, deletion of an AzureKeyVaultSecret resource
func (c *Controller) handleKvSecretChange(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.Errorf("Failed to get namespace and name from %s: %+v", key, err)
		return err
	}

	secretClient := c.kubeclientset.CoreV1().Secrets(namespace)

	kvSecretResource, err := c.kvSecretsLister.AzureKeyVaultSecrets(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.Infof("Deleting secret \"%s\"", key)
		secretClient.Delete(name, v1.NewDeleteOptions(60))
		return nil
	} else if err != nil {
		klog.Errorf("Failed to fetch AzureKeyVaultSecret resource with name %s: %+v", key, err)
		return err
	}

	url := kvSecretResource.Spec.KeyVaultId
	data := make(map[string][]byte)

	for _, secret := range kvSecretResource.Spec.SecretList {
		value, err := c.keyVaultClient.GetSecret(url, secret.SecretKey)
		if err != nil {
			klog.Warningf("Could not get Azure Key Vault secret \"%s\" from %s: %+v ...ignoring", secret.SecretKey, url, err)
			continue
		}

		data[secret.TargetSecretName] = []byte(value)
	}

	secret, err := secretClient.Get(name, v1.GetOptions{})
	if err != nil {
		klog.Infof("Creating new secret \"%s\"", key)
		_, err = secretClient.Create(
			&corev1.Secret{
				Type: corev1.SecretTypeOpaque,
				ObjectMeta: v1.ObjectMeta{
					Name: name,
				},
				Data: data,
			},
		)

		if err != nil {
			klog.Errorf("Failed to create secret \"%s\": %+v", key, err)
		}
	} else {
		klog.Infof("Updating existing secret \"%s\"", key)
		secret.Data = data
		_, err = secretClient.Update(secret)

		if err != nil {
			klog.Errorf("Failed to update secret \"%s\": %+v", key, err)
		}
	}

	return nil
}
