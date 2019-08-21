package kitaspace

import (
	"math/rand"
	"os"

	kitav1alpha1 "github.com/danistrebel/kita-operator/pkg/apis/kita/v1alpha1"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("kitaspace_auth_token")

const tokenKey string = "token"

func newLoginTokenForCR(ksr *kitav1alpha1.KitaSpace, scheme *runtime.Scheme) (*corev1.Secret, error) {
	labels := map[string]string{
		"app":   ksr.Name,
		"owner": ksr.Spec.Owner.Name,
	}

	generatedToken := generateRandomToken(16)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tokenSecretName(ksr),
			Namespace: ksr.Name,
			Labels:    labels,
		},
		StringData: map[string]string{
			"token": generatedToken,
		},
	}

	// Set KitaSpace instance as the owner and controller
	if err := controllerutil.SetControllerReference(ksr, secret, scheme); err != nil {
		return secret, err
	}

	return secret, nil
}

func newServiceAccountForCR(ksr *kitav1alpha1.KitaSpace, scheme *runtime.Scheme) (*corev1.ServiceAccount, error) {
	labels := map[string]string{
		"app": ksr.Name,
	}

	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "space-admin",
			Namespace: ksr.Name,
			Labels:    labels,
		},
	}

	// Set KitaSpace instance as the owner and controller
	if err := controllerutil.SetControllerReference(ksr, sa, scheme); err != nil {
		return sa, err
	}

	return sa, nil
}

func newRoleBindingForCR(ksr *kitav1alpha1.KitaSpace, scheme *runtime.Scheme) (*rbacv1.RoleBinding, error) {

	rb := &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind: "RoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "space-admin-binding",
			Namespace: ksr.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "space-admin",
				Namespace: ksr.Name,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "admin",
		},
	}

	// Set KitaSpace instance as the owner and controller
	if err := controllerutil.SetControllerReference(ksr, rb, scheme); err != nil {
		return rb, err
	}

	return rb, nil
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

func sendTokenAsEmail(token string, ksr *kitav1alpha1.KitaSpace) {

	sendGridAPIKey, exists := os.LookupEnv("SENDGRID_API_KEY")

	if !exists {
		log.Info("No Sendgrid configured. Cannot send token for space: " + ksr.Name)
		return
	}

	from := mail.NewEmail("Kita Operator", os.Getenv("SENDGRID_EMAIL_SENDER"))
	subject := "Your Kita Password"
	to := mail.NewEmail(ksr.Spec.Owner.Name, ksr.Spec.Owner.Email)
	plainTextContent := "Your Kita Token is: " + token
	htmlContent := "Your Kita Token is: <strong>" + token + "</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendGridAPIKey)
	result, err := client.Send(message)

	if err != nil {
		log.Error(err, "Error while sending Email")
	} else {
		log.Info("Sent token as email to: " + ksr.Spec.Owner.Email)
		log.Info(result.Body)
	}
}
