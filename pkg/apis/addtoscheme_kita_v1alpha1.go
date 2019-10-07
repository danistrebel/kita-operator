package apis

import (
	"github.com/danistrebel/kita-operator/pkg/apis/kita/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, routev1.SchemeBuilder.AddToScheme)
}
