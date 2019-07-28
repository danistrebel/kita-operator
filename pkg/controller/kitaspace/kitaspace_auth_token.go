package kitaspace

import (
	"math/rand"

	kitav1alpha1 "github.com/danistrebel/kita-operator/pkg/apis/kita/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const tokenKey string = "token"

func newLoginTokenForCR(ksr *kitav1alpha1.KitaSpace, scheme *runtime.Scheme) (*corev1.Secret, error) {
	labels := map[string]string{
		"app":   ksr.Name,
		"owner": ksr.Spec.Owner.Name,
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tokenSecretName(ksr),
			Namespace: ksr.Namespace,
			Labels:    labels,
		},
		StringData: map[string]string{
			"token": generateRandomToken(16),
		},
	}

	// Set KitaSpace instance as the owner and controller
	if err := controllerutil.SetControllerReference(ksr, secret, scheme); err != nil {
		return secret, err
	}

	return secret, nil
}

func tokenSecretName(ksr *kitav1alpha1.KitaSpace) string {
	return ksr.Name + "-token"
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generate a random string token of length n
func generateRandomToken(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
