package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AzureKeyVaultSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureKeyVaultSecretSpec   `json:"spec"`
	Status AzureKeyVaultSecretStatus `json:"status"`
}

type AzureKeyVaultSecretSpec struct {
	KeyVaultId string      `json:"keyVaultId"`
	SecretList []SecretMap `json:"secretList"`
}

type SecretMap struct {
	SecretKey        string `json:"key"`
	TargetSecretName string `json:"targetSecretName"`
}

type AzureKeyVaultSecretStatus struct {
	State string `json:"state"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AzureKeyVaultSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AzureKeyVaultSecret `json:"items"`
}
