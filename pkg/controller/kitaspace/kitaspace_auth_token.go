package kitaspace

import (
	"math/rand"
	"os"

	kitav1alpha1 "github.com/danistrebel/kita-operator/pkg/apis/kita/v1alpha1"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	corev1 "k8s.io/api/core/v1"
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
			Namespace: ksr.Namespace,
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
	from := mail.NewEmail("Kita Operator", os.Getenv("SENDGRID_EMAIL_SENDER"))
	subject := "Your Kita Password"
	to := mail.NewEmail(ksr.Spec.Owner.Name, ksr.Spec.Owner.Email)
	plainTextContent := "Your Kita Token is: " + token
	htmlContent := "Your Kita Token is: <strong>" + token + "</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	result, err := client.Send(message)

	if err != nil {
		log.Error(err, "Error while sending Email")
	} else {
		log.Info("Sent token as email to: " + ksr.Spec.Owner.Email)
		log.Info(result.Body)
	}
}