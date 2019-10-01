package kitaspace

import (
	"strconv"

	kitav1alpha1 "github.com/danistrebel/kita-operator/pkg/apis/kita/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilpointer "k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const terminalPortWebPortName = "web"
const terminalPortWebPort = 9000

func defaultLabels(cr *kitav1alpha1.KitaSpace) map[string]string {
	return map[string]string{
		"app":   cr.Name,
		"owner": cr.Spec.Owner.Name,
	}
}

// newKitaTerminalDeloymentforCR returns a kita terminal deployment in the appropriate namespace
func newKitaTerminalDeloymentforCR(cr *kitav1alpha1.KitaSpace, scheme *runtime.Scheme) (*appsv1.Deployment, error) {

	initContainers := getRepoInitContainers(cr.Spec.Repos)

	vsCodeImage := "chinodesuuu/coder:vanilla"
	if cr.Spec.Platform == "openshift" {
		vsCodeImage = "chinodesuuu/coder:openshift"
	}
	print(vsCodeImage)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cli-terminal-deployment",
			Namespace: cr.Name,
			Labels:    defaultLabels(cr),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utilpointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": cr.Name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cli-terminal",
					Namespace: cr.Name,
					Labels:    defaultLabels(cr),
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "space-admin",
					InitContainers:     initContainers,
					Containers: []corev1.Container{
						{
							Name:  "vscode",
							Image: vsCodeImage,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: terminalPortWebPort,
									Name:          terminalPortWebPortName,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "CODER_ENABLE_AUTH",
									Value: "true",
								},
								{
									Name: "CODER_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: v1.LocalObjectReference{
												Name: tokenSecretName(cr),
											},
											Key: tokenKey,
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "project-files",
									MountPath: "/home/coder/projects",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "project-files",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	// Set KitaSpace instance as the owner and controller
	if err := controllerutil.SetControllerReference(cr, deployment, scheme); err != nil {
		return deployment, err
	}

	return deployment, nil
}

// Create init containers for got repositories
func getRepoInitContainers(repos []string) []corev1.Container {
	var repoInitContainers []corev1.Container
	for i, repo := range repos {
		initContainer := corev1.Container{
			Name:  "git-init-" + strconv.Itoa(i),
			Image: "alpine/git",
			Command: []string{
				"git",
				"clone",
				repo,
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "project-files",
					MountPath: "/git",
				},
			},
		}

		repoInitContainers = append(repoInitContainers, initContainer)
	}

	return repoInitContainers
}

func newKitaTerminalServiceForCR(cr *kitav1alpha1.KitaSpace, scheme *runtime.Scheme) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cli-terminal",
			Namespace: cr.Name,
			Labels:    defaultLabels(cr),
		},
		Spec: corev1.ServiceSpec{
			Selector: defaultLabels(cr),
			Ports: []corev1.ServicePort{
				{
					Port: terminalPortWebPort,
					Name: terminalPortWebPortName,
				},
			},
		},
	}

	// Set KitaSpace instance as the owner and controller
	if err := controllerutil.SetControllerReference(cr, service, scheme); err != nil {
		return service, err
	}

	return service, nil
}
