// +build !

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1.KitaSpace":       schema_pkg_apis_kita_v1alpha1_KitaSpace(ref),
		"github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1.KitaSpaceSpec":   schema_pkg_apis_kita_v1alpha1_KitaSpaceSpec(ref),
		"github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1.KitaSpaceStatus": schema_pkg_apis_kita_v1alpha1_KitaSpaceStatus(ref),
		"github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1.OwnerSpec":       schema_pkg_apis_kita_v1alpha1_OwnerSpec(ref),
	}
}

func schema_pkg_apis_kita_v1alpha1_KitaSpace(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KitaSpace is the Schema for the kitaspaces API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1.KitaSpaceSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1.KitaSpaceStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1.KitaSpaceSpec", "github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1.KitaSpaceStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_kita_v1alpha1_KitaSpaceSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KitaSpaceSpec defines the desired state of KitaSpace",
				Properties: map[string]spec.Schema{
					"owner": {
						SchemaProps: spec.SchemaProps{
							Description: "INSERT ADDITIONAL SPEC FIELDS - desired state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html",
							Ref:         ref("github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1.OwnerSpec"),
						},
					},
					"repos": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
				},
				Required: []string{"owner"},
			},
		},
		Dependencies: []string{
			"github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1.OwnerSpec"},
	}
}

func schema_pkg_apis_kita_v1alpha1_KitaSpaceStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KitaSpaceStatus defines the observed state of KitaSpace",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_kita_v1alpha1_OwnerSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OwnerSpec defines the desired state of KitaSpace",
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"email": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"name", "email"},
			},
		},
		Dependencies: []string{},
	}
}
