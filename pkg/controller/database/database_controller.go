package database

import (
	"context"
	dbv1alpha1 "db-operator/pkg/apis/db/v1alpha1"
	"fmt"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
	"time"
)

var log = logf.Log.WithName("controller_database")

const dbFinalizer = "finalizer.db.clarizen.cloud"

// Add creates a new Database Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDatabase{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("database-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	pred := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			// Metadata.generation changes if Spec was changed
			log.Info("Update event")

			return inWatchedNamespace(e.MetaNew.GetNamespace()) && e.MetaOld.GetGeneration() != e.MetaNew.GetGeneration()
		},
		CreateFunc: func(e event.CreateEvent) bool {
			watchTime := time.Now().Add(-1 * time.Minute)
			creationTime := e.Meta.GetCreationTimestamp()

			return inWatchedNamespace(e.Meta.GetNamespace()) && creationTime.After(watchTime)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return inWatchedNamespace(e.Meta.GetNamespace())
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return inWatchedNamespace(e.Meta.GetNamespace())
		},
	}

	// Watch for changes to primary resource Database
	err = c.Watch(&source.Kind{Type: &dbv1alpha1.Database{}}, &handler.EnqueueRequestForObject{}, pred)
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileDatabase implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDatabase{}

// ReconcileDatabase reconciles a Database object
type ReconcileDatabase struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Database object and makes changes based on the state read
// and what is in the Database.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDatabase) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Database")

	// Fetch the Database instance
	instance := &dbv1alpha1.Database{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	isDbMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isDbMarkedToBeDeleted {
		if contains(instance.GetFinalizers(), dbFinalizer) {
			// Run finalization logic for dbFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeDatabase(reqLogger, instance); err != nil {
				return reconcile.Result{}, err
			}

			// Remove dbFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			instance.SetFinalizers(remove(instance.GetFinalizers(), dbFinalizer))
			err := r.client.Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// Add finalizer for this CR
	if !contains(instance.GetFinalizers(), dbFinalizer) {
		if err := r.addFinalizer(reqLogger, instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	usr := &user{}
	// Check if this Database already exists and status is "Created"
	err = updateEvent(instance, usr)
	if err != nil {
		return reconcile.Result{}, err
	}

	// If secret was set, we have to create k8s secret
	if usr.password != "" {
		secret, err := updateSecret(instance, usr)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = r.client.Create(context.TODO(), secret)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Database created/updated successfully - don't requeue
	return reconcile.Result{}, nil
}

func (r *ReconcileDatabase) finalizeDatabase(reqLogger logr.Logger, m *dbv1alpha1.Database) error {
	err := deleteEvent(m)
	if err != nil {
		return err
	}

	err = r.client.Delete(context.TODO(), &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-db-secret", m.Name),
			Namespace: m.Namespace,
		},
	})

	reqLogger.Info("Successfully finalized database")
	return nil
}

func (r *ReconcileDatabase) addFinalizer(reqLogger logr.Logger, m *dbv1alpha1.Database) error {
	reqLogger.Info("Adding Finalizer for the Database")
	m.SetFinalizers(append(m.GetFinalizers(), dbFinalizer))

	// Update CR
	err := r.client.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update Database with finalizer")
		return err
	}
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func remove(list []string, s string) []string {
	for i, v := range list {
		if v == s {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func inWatchedNamespace(ns string) bool {
	watchedNamespaces := strings.Split(os.Getenv("NAMESPACES"), ",")

	return contains(watchedNamespaces, ns)
}
