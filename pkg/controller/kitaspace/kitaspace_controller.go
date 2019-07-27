package kitaspace

import (
	"context"
	"strconv"

	kitav1alpha1 "github.com/danistrebel/kita-space-operator/pkg/apis/kita/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_kitaspace")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new KitaSpace Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKitaSpace{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("kitaspace-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KitaSpace
	err = c.Watch(&source.Kind{Type: &kitav1alpha1.KitaSpace{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner KitaSpace
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kitav1alpha1.KitaSpace{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKitaSpace implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKitaSpace{}

// ReconcileKitaSpace reconciles a KitaSpace object
type ReconcileKitaSpace struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a KitaSpace object and makes changes based on the state read
// and what is in the KitaSpace.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKitaSpace) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KitaSpace")

	// Fetch the KitaSpace instance
	instance := &kitav1alpha1.KitaSpace{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set KitaSpace instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}
		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

// Create init containers for repos
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

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *kitav1alpha1.KitaSpace) *corev1.Pod {
	labels := map[string]string{
		"app":   cr.Name,
		"owner": cr.Spec.Owner.Name,
	}

	initContainers := getRepoInitContainers(cr.Spec.Repos)

	return &corev1.Pod{
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
							Name:  "PASSWORD",
							Value: "replaceme",
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
}
