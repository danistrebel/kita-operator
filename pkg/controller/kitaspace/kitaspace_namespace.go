package kitaspace

import (
	kitav1alpha1 "github.com/danistrebel/kita-operator/pkg/apis/kita/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func newNamespaceForCR(ksr *kitav1alpha1.KitaSpace, scheme *runtime.Scheme) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ksr.Name,
		},
	}

	// Set KitaSpace instance as the owner and controller
	if err := controllerutil.SetControllerReference(ksr, namespace, scheme); err != nil {
		return namespace, err
	}

	return namespace, nil
}
