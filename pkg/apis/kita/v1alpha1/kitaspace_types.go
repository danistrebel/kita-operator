package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KitaSpaceSpec defines the desired state of KitaSpace
// +k8s:openapi-gen=true
type KitaSpaceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Owner    OwnerSpec `json:"owner"`
	Repos    []string  `json:"repos,omitempty"`    // OPTIONAL: Git Repos to be initialized in the workspace
	Platform string    `json:"platform,omitempty"` // OPTIONAL: Set to "openshift" for Openshift configuration
	Token    string    `json:"token,omitempty"`    // OPTIONAL: Explicitly set login token for Coder
}

// OwnerSpec defines the desired state of KitaSpace
// +k8s:openapi-gen=true
type OwnerSpec struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// KitaSpaceStatus defines the observed state of KitaSpace
// +k8s:openapi-gen=true
type KitaSpaceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KitaSpace is the Schema for the kitaspaces API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type KitaSpace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KitaSpaceSpec   `json:"spec,omitempty"`
	Status KitaSpaceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KitaSpaceList contains a list of KitaSpace
type KitaSpaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KitaSpace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KitaSpace{}, &KitaSpaceList{})
}
