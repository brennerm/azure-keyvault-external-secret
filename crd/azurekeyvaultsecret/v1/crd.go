package v1

import (
	"reflect"

	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

const (
	CRDPlural   string = "azurekeyvaultsecrets"
	CRDGroup    string = "brennerm.github.io"
	CRDVersion  string = "v1"
	FullCRDName string = CRDPlural + "." + CRDGroup
)

func CreateCRD(clientset clientset.Interface) error {
	klog.Info("Creating CRDs")

	crd := &apiextensionv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: FullCRDName},

		Spec: apiextensionv1.CustomResourceDefinitionSpec{
			Group: CRDGroup,
			Scope: apiextensionv1.NamespaceScoped,
			Versions: []apiextensionv1.CustomResourceDefinitionVersion{
				apiextensionv1.CustomResourceDefinitionVersion{
					Name:    CRDVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextensionv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionv1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensionv1.JSONSchemaProps{
								"spec": apiextensionv1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]apiextensionv1.JSONSchemaProps{
										"keyVaultId": apiextensionv1.JSONSchemaProps{
											Type: "string",
										},
										"secretList": apiextensionv1.JSONSchemaProps{
											Type: "array",
											Items: &apiextensionv1.JSONSchemaPropsOrArray{
												Schema: &apiextensionv1.JSONSchemaProps{
													Type: "object",
													Properties: map[string]apiextensionv1.JSONSchemaProps{
														"key": apiextensionv1.JSONSchemaProps{
															Type: "string",
														},
														"targetSecretName": apiextensionv1.JSONSchemaProps{
															Type: "string",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
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
		klog.Info("CRDs have been created")
		return nil
	}
	return err
}
