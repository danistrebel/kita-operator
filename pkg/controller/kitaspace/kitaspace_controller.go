package kitaspace

import (
	"context"

	kitav1alpha1 "github.com/danistrebel/kita-operator/pkg/apis/kita/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

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
	c, err := controller.New("kitaspace-controller", mgr, controller.Options{
		Reconciler: r,
	})
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
//
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKitaSpace) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KitaSpace")

	///////////////////////
	// KITA SPACE CR
	///////////////////////
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

	///////////////////////
	// KITA SPACE NAMESPACE
	///////////////////////
	//
	// Create a new dedicated namespace the kita space
	// Hackish because r.client.Get apparently does not support namespaces
	namespace, err := newNamespaceForCR(instance, r.scheme)
	if err != nil {
		return reconcile.Result{}, err
	}
	err = r.client.Create(context.TODO(), namespace)
	if err != nil && !errors.IsAlreadyExists(err) {
		reqLogger.Info(err.Error())
		return reconcile.Result{}, err
	} else if err == nil {
		return reconcile.Result{Requeue: true}, nil
	} else {
		reqLogger.Info("Namespace Exists")
	}

	///////////////////////
	// KITA SPACE TOKEN
	///////////////////////
	loginToken, err := newLoginTokenForCR(instance, r.scheme)
	if err != nil {
		return reconcile.Result{}, err
	}
	foundLoginTokenSecret := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: loginToken.Name, Namespace: instance.Name}, foundLoginTokenSecret)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new login Secret", "Secret.Namespace", loginToken.Namespace, "Secret.Name", loginToken.Name)
		err = r.client.Create(context.TODO(), loginToken)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Send Email
		sendTokenAsEmail(loginToken.StringData[tokenKey], instance)

		// Secret created successfully - Requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	///////////////////////
	// KITA SERVICE ACCOUNT
	///////////////////////
	sa, err := newServiceAccountForCR(instance, r.scheme)
	if err != nil {
		return reconcile.Result{}, err
	}
	foundServiceAccount := &corev1.ServiceAccount{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: sa.Name, Namespace: instance.Name}, foundServiceAccount)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service Account", "sa.Namespace", sa.Namespace, "sa.Name", sa.Name)
		err = r.client.Create(context.TODO(), sa)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Secret created successfully - Requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	///////////////////////
	// KITA ROLE BINDING
	///////////////////////
	rb, err := newRoleBindingForCR(instance, r.scheme)
	if err != nil {
		return reconcile.Result{}, err
	}
	foundRoleBinding := &rbacv1.RoleBinding{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: rb.Name, Namespace: instance.Name}, foundRoleBinding)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Role Binding", "sa.Namespace", rb.Namespace, "sa.Name", rb.Name)
		err = r.client.Create(context.TODO(), rb)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Role Binding created successfully - Requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	///////////////////////
	// KITA SPACE TERMINAL DEPLOYMENT
	///////////////////////
	// Define a new Deployment object
	deployment, err := newKitaTerminalDeloymentforCR(instance, r.scheme)
	// Set KitaSpace instance as the owner and controller
	if err != nil {
		return reconcile.Result{}, err
	}
	// Check if this Deployment already exists
	foundDeployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: instance.Name}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "deployment.Namespace", instance.Name, "deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}
		// Deployment created successfully - requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	///////////////////////
	// KITA SPACE TERMINAL SERVICE
	///////////////////////
	// Define a new Service object
	service, err := newKitaTerminalServiceForCR(instance, r.scheme)
	// Check if this Service already exists
	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: instance.Name}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "service.Namespace", instance.Name, "service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}
		// Service created successfully - requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	///////////////////////
	// KITA SPACE TERMINAL OPENSHIFT ROUTE
	///////////////////////

	if instance.Spec.Platform == "openshift" {
		// Define a new Route object
		route, err := newKitaTerminalRouteForCR(instance, r.scheme)
		// Check if this Deployment already exists
		foundRoute := &routev1.Route{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: route.Name, Namespace: instance.Name}, foundRoute)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Route", "route.Namespace", instance.Name, "route.Name", route.Name)
			err = r.client.Create(context.TODO(), route)
			if err != nil {
				return reconcile.Result{}, err
			}
			// Route created successfully - requeue
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			return reconcile.Result{}, err
		}
	}

	// Resources already exist - don't requeue
	reqLogger.Info("Skip reconcile: All KitaSpace resources already exist")

	return reconcile.Result{}, nil
}
