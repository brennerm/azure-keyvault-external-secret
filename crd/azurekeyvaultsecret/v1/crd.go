package v1

import (
	"reflect"

	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CRDPlural   string = "secrets"
	CRDGroup    string = "azureKeyvault"
	CRDVersion  string = "v1alpha1"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

func CreateCRD(namespace string, clientset apiextension.Interface) error {
	crd := &apiextensionv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      FullCRDName,
			Namespace: namespace,
		},
		Spec: apiextensionv1.CustomResourceDefinitionSpec{
			Group: CRDGroup,
			Scope: apiextensionv1.NamespaceScoped,
			Versions: []apiextensionv1.CustomResourceDefinitionVersion{
				apiextensionv1.CustomResourceDefinitionVersion{
					Name:   CRDVersion,
					Served: true,
				},
			},
			Names: apiextensionv1.CustomResourceDefinitionNames{
				Plural: CRDPlural,
				Kind:   reflect.TypeOf(AzureKeyVaultSecret{}).Name(),
			},
		},
	}

	_, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Create(crd)
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}
