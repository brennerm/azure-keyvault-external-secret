package main

import (
	"flag"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"k8s.io/sample-controller/pkg/signals"

	v1 "github.com/brennerm/azure-keyvault-external-secret/crd/azurekeyvaultsecret/v1"
	clientset "github.com/brennerm/azure-keyvault-external-secret/crd/generated/clientset/versioned"
	informers "github.com/brennerm/azure-keyvault-external-secret/crd/generated/informers/externalversions"
	"github.com/brennerm/azure-keyvault-external-secret/keyvault"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	exampleClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	extensionClient, err := apiextension.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Failed to create client: %v", err)
	}

	keyVaultClient, err := keyvault.NewKeyVaultClient()
	if err != nil {
		klog.Errorf("Failed to create Key Vault Client: %+v", err)
	}

	// Create the CRD
	err = v1.CreateCRD(extensionClient)
	if err != nil {
		klog.Fatalf("Failed to create crd: %v", err)
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*60)
	akesInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*60)

	controller := NewController(
		kubeClient,
		exampleClient,
		kubeInformerFactory.Core().V1().Secrets(),
		akesInformerFactory.Azurekeyvaultsecret().V1().AzureKeyVaultSecrets(),
		keyVaultClient,
	)

	kubeInformerFactory.Start(stopCh)
	akesInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
