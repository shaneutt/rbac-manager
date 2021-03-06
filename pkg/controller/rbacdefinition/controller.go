/*
Copyright 2018 ReactiveOps.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rbacdefinition

import (
	"context"

	rbacmanagerv1beta1 "github.com/reactiveops/rbac-manager/pkg/apis/rbacmanager/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new RBACDefinition Controller and adds it to the Manager.
// The Manager will set fields on the Controller and Start it.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())

	if err != nil {
		// If we can't get a clientset we can't do anything else
		panic(err)
	}

	return &ReconcileRBACDefinition{
		Client:    mgr.GetClient(),
		clientset: clientset,
		scheme:    mgr.GetScheme(),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("rbacdefinition-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to RBACDefinition
	err = c.Watch(&source.Kind{
		Type: &rbacmanagerv1beta1.RBACDefinition{},
	}, &handler.EnqueueRequestForObject{})

	if err != nil {
		return err
	}

	return nil
}

// ReconcileRBACDefinition reconciles a RBACDefinition object
type ReconcileRBACDefinition struct {
	client.Client
	scheme    *runtime.Scheme
	clientset kubernetes.Interface
}

// Reconcile makes changes in response to RBACDefinition changes
func (r *ReconcileRBACDefinition) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	var err error
	rdr := Reconciler{Clientset: r.clientset}

	// Fetch the RBACDefinition instance
	rbacDef := &rbacmanagerv1beta1.RBACDefinition{}
	err = r.Get(context.TODO(), request.NamespacedName, rbacDef)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	rdr.Reconcile(rbacDef)

	return reconcile.Result{}, nil
}
