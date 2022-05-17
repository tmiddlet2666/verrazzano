// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package configmaps

import (
	"context"
	os2 "github.com/verrazzano/verrazzano/pkg/os"
	"github.com/verrazzano/verrazzano/platform-operator/constants"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	vzctrl "github.com/verrazzano/verrazzano/pkg/controller"
	"github.com/verrazzano/verrazzano/pkg/log/vzlog"
	installv1alpha1 "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/controllers"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"

	"go.uber.org/zap"
)

// VerrazzanoConfigMapsReconciler reconciles ConfigMaps.
// This controller manages install override sources from the Verrazzano CR
type VerrazzanoConfigMapsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	log    vzlog.VerrazzanoLogger
}

// SetupWithManager creates a new controller and adds it to the manager
func (r *VerrazzanoConfigMapsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Complete(r)
}

// Reconcile the ConfigMap
func (r *VerrazzanoConfigMapsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	if ctx == nil {
		ctx = context.TODO()
	}

	// Get Verrazzano from the cluster
	vzList := &installv1alpha1.VerrazzanoList{}
	err := r.List(ctx, vzList)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		zap.S().Errorf("Failed to fetch Verrazzano resource: %v", err)
		return newRequeueWithDelay(), err
	}

	if vzList != nil && len(vzList.Items) > 0 {
		vz := &vzList.Items[0]
		res, err := r.reconcileInstallOverrideConfigMap(ctx, req, vz)
		if err != nil {
			zap.S().Errorf("Failed to reconcile ConfigMap: %v", err)
			return newRequeueWithDelay(), err
		}
		return res, nil
	}
	return ctrl.Result{}, nil
}

// reconcileInstallOverrideConfigMap looks through the Verrazzano CR for the ConfigMap
// if the request is from the same namespace as the CR
func (r *VerrazzanoConfigMapsReconciler) reconcileInstallOverrideConfigMap(ctx context.Context, req ctrl.Request, vz *installv1alpha1.Verrazzano) (ctrl.Result, error) {

	// Get the ConfigMap present in the Verrazzano CR namespace
	configMap := &corev1.ConfigMap{}
	if vz.Namespace == req.Namespace {
		if err := r.Get(ctx, req.NamespacedName, configMap); err != nil {
			// Do not reconcile if the ConfigMap was deleted
			if errors.IsNotFound(err) {
				if err := controllers.ProcDeletedOverride(r.Client, vz, req.Name, constants.ConfigMapKind); err != nil {
					// Do not return an error as it's most likely due to timing
					return newRequeueWithDelay(), nil
				}
				return reconcile.Result{}, nil
			}
			zap.S().Errorf("Failed to fetch ConfigMap in Verrazzano CR namespace: %v", err)
			return newRequeueWithDelay(), err
		}

		if result, err := r.initLogger(*configMap); err != nil {
			return result, err
		}

		componentCtx, err := spi.NewContext(r.log, r.Client, vz, false)
		if err != nil {
			r.log.Errorf("Failed to construct component context: %v", err)
			return newRequeueWithDelay(), err
		}

		// Check if the ConfigMap is listed as an override source under a component
		if componentName, ok := controllers.VzContainsResource(componentCtx, configMap.Name, configMap.Kind); ok {
			if configMap.DeletionTimestamp.IsZero() {
				// Check if our finalizer is already present
				if !controllerutil.ContainsFinalizer(configMap, constants.OverridesFinalizer) {
					configMap.Finalizers = append(configMap.Finalizers, constants.OverridesFinalizer)
					err := r.Update(context.TODO(), configMap)
					if err != nil {
						return newRequeueWithDelay(), err
					}
					return reconcile.Result{Requeue: true}, nil
				}
			} else {
				// Requeue if other finalizers are present
				if configMap.Finalizers != nil && !controllerutil.ContainsFinalizer(configMap, constants.OverridesFinalizer) {
					return reconcile.Result{Requeue: true}, nil
				}

				// Now since only our finalizer is present, therefore we remove it to delete the ConfigMap
				// and trigger verrazzano reconcile
				controllerutil.RemoveFinalizer(configMap, constants.OverridesFinalizer)
				err := r.Update(context.TODO(), configMap)
				if err != nil {
					return newRequeueWithDelay(), err
				}
			}

			err := controllers.UpdateVerrazzanoForInstallOverrides(r.Client, componentCtx, componentName)
			if err != nil {
				r.log.ErrorfThrottled("Failed to reconcile ConfigMap: %v", err)
				return newRequeueWithDelay(), err
			}
			r.log.Infof("Updated Verrazzano Resource")
		}
	}
	return ctrl.Result{}, nil
}

// GetInstallOverridesFromConfigMap takes the list of Overrides and returns a list of key value pairs
func GetInstallOverridesFromConfigMap(ctx spi.ComponentContext, overrides []installv1alpha1.Overrides) ([]string, error) {
	var file *os.File
	var files []string

	for _, override := range overrides {
		// Check if ConfigMapRef is populated and gather helm file
		if override.ConfigMapRef != nil {
			// Get the ConfigMap
			configMap := &corev1.ConfigMap{}
			selector := override.ConfigMapRef
			nsn := types.NamespacedName{Name: selector.Name, Namespace: ctx.EffectiveCR().Namespace}
			optional := selector.Optional
			err := ctx.Client().Get(context.TODO(), nsn, configMap)
			if err != nil {
				if optional == nil || !*optional {
					err := ctx.Log().ErrorfNewErr("Could not get Configmap %s from namespace %s: %v", nsn.Name, nsn.Namespace, err)
					return files, err
				}
				ctx.Log().Debugf("Optional Configmap %s from namespace %s not found", nsn.Name, nsn.Namespace)
				continue
			}

			// Get resource data
			fieldData, ok := configMap.Data[selector.Key]
			if !ok {
				if optional == nil || !*optional {
					err := ctx.Log().ErrorfNewErr("Could not get Data field %s from Resource %s from namespace %s", selector.Key, nsn.Name, nsn.Namespace)
					return files, err
				}
				ctx.Log().Debugf("Optional Resource %s from namespace %s missing Data key %s", nsn.Name, nsn.Namespace, selector.Key)
			}

			// Create the temp file for the data
			file, err = os2.CreateTempFile(ctx.Log(), "install-overrides-*.yaml", []byte(fieldData))
			if err != nil {
				return files, err
			}
			files = append(files, file.Name())
		}
	}
	return files, nil
}

// initialize logger for ConfigMap
func (r *VerrazzanoConfigMapsReconciler) initLogger(cm corev1.ConfigMap) (ctrl.Result, error) {
	// Get the resource logger needed to log message using 'progress' and 'once' methods
	log, err := vzlog.EnsureResourceLogger(&vzlog.ResourceConfig{
		Name:           cm.Name,
		Namespace:      cm.Namespace,
		ID:             string(cm.UID),
		Generation:     cm.Generation,
		ControllerName: "ConfigMaps",
	})
	if err != nil {
		zap.S().Errorf("Failed to create resource logger for VerrazzanoConfigMap controller: %v", err)
		return newRequeueWithDelay(), err
	}
	r.log = log
	return ctrl.Result{}, nil
}

// Create a new Result that will cause a reconcile requeue after a short delay
func newRequeueWithDelay() ctrl.Result {
	return vzctrl.NewRequeueWithDelay(3, 5, time.Second)
}
