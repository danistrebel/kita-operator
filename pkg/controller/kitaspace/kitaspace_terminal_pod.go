package kitaspace

import (
	"strconv"

	kitav1alpha1 "github.com/danistrebel/kita-operator/pkg/apis/kita/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// newKitaTerminalPodForCR returns a kita terminal pod with the same name/namespace as the cr
func newKitaTerminalPodForCR(cr *kitav1alpha1.KitaSpace, scheme *runtime.Scheme) (*corev1.Pod, error) {
	labels := map[string]string{
		"app":   cr.Name,
		"owner": cr.Spec.Owner.Name,
	}

	initContainers := getRepoInitContainers(cr.Spec.Repos)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Owner.Name + "-terminal",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			InitContainers: initContainers,
			Containers: []corev1.Container{
				{
					Name:  "vscode",
					Image: "codercom/code-server",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 8443,
							Name:          "web",
						},
					},
					Env: []corev1.EnvVar{
						{
							Name: "PASSWORD",
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
							MountPath: "/home/coder/project",
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
	}

	// Set KitaSpace instance as the owner and controller
	if err := controllerutil.SetControllerReference(cr, pod, scheme); err != nil {
		return pod, err
	}

	return pod, nil
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
