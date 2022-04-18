// Copyright (c) 2020, 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package verrazzano

import (
	"context"
	"github.com/verrazzano/verrazzano/pkg/log/vzlog"
	installv1alpha1 "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Watch configmaps that hold helm values
// The reconciler will be called if these are referenced in the Verrazzano CR
func (r *Reconciler) watchConfigMaps(namespace string, name string, log vzlog.VerrazzanoLogger) error {
	// Define a mapping to the Verrazzano resource
	mapFn := handler.EnqueueRequestsFromMapFunc(
		func(a client.Object) []reconcile.Request {
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: namespace,
						Name:      name,
					},
				},
			}
		})

	// Get the Verrazzano Resource
	vz := &installv1alpha1.Verrazzano{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, vz)
	if err != nil {
		err = log.ErrorfNewErr("Could not get the Verrazzano resource %s/%s, error: %v", namespace, name, err)
		return err
	}

	// Watch ConfigMap create
	predicateFunc := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			// Cast object to ConfigMap
			configmap := e.Object.(*corev1.ConfigMap)

			// Filter events only in the Verrazzano namespace
			if configmap.Namespace != namespace {
				return false
			}

			// Verify that Verrazzano contains the given resource
			if !vzContainsResource(vz, e.Object) {
				return false
			}

			log.Debugf("Configmap %s in namespace %s found in the Verrazzano CR", configmap.Name, configmap.Namespace)
			return true
		},
	}

	// Watch ConfigMaps and trigger reconciles for Verrazzano resources when a ConfigMap is updated with the correct criteria
	err = r.Controller.Watch(
		&source.Kind{Type: &corev1.ConfigMap{}},
		mapFn,
		predicateFunc)
	if err != nil {
		return err
	}
	log.Debugf("Watching for Configmaps to activate reconcile for Verrazzano CR %s/%s", namespace, name)
	return nil
}

// Watch configmaps that hold helm values
// The reconciler will be called if these are referenced in the Verrazzano CR
func (r *Reconciler) watchSecrets(namespace string, name string, log vzlog.VerrazzanoLogger) error {
	// Define a mapping to the Verrazzano resource
	mapFn := handler.EnqueueRequestsFromMapFunc(
		func(a client.Object) []reconcile.Request {
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: namespace,
						Name:      name,
					},
				},
			}
		})

	// Get the Verrazzano Resource
	vz := &installv1alpha1.Verrazzano{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, vz)
	if err != nil {
		err = log.ErrorfNewErr("Could not get the Verrazzano resource %s/%s, error: %v", namespace, name, err)
		return err
	}

	// Watch ConfigMap create
	predicateFunc := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			// Cast object to Secret
			secret := e.Object.(*corev1.Secret)

			// Filter events only in the Verrazzano namespace
			if secret.Namespace != namespace {
				return false
			}

			// Verify that Verrazzano contains the given resource
			if !vzContainsResource(vz, e.Object) {
				return false
			}

			log.Debugf("Secret %s in namespace %s found in the Verrazzano CR", secret.Name, secret.Namespace)
			return true
		},
	}

	// Watch ConfigMaps and trigger reconciles for Verrazzano resources when a ConfigMap is updated with the correct criteria
	err = r.Controller.Watch(
		&source.Kind{Type: &corev1.Secret{}},
		mapFn,
		predicateFunc)
	if err != nil {
		return err
	}
	log.Debugf("Watching for Secrets to activate reconcile for Verrazzano CR %s/%s", namespace, name)
	return nil
}

// vzContainsResource checks to see if the resource is listed in the Verrazzano
func vzContainsResource(vz *installv1alpha1.Verrazzano, object client.Object) bool {
	// TODO: verify that the Verrazzano contains the object in a component helm config
	return true
}
