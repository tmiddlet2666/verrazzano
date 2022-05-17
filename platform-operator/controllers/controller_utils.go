// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package controllers

import (
	"context"
	"fmt"
	"github.com/verrazzano/verrazzano/pkg/log/vzlog"
	installv1alpha1 "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/constants"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/registry"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// VzContainsResource checks to see if the resource is listed in the Verrazzano
func VzContainsResource(ctx spi.ComponentContext, objectName string, objectKind string) (string, bool) {
	for _, component := range registry.GetComponents() {
		if component.MonitorOverrides(ctx) {
			if found := componentContainsResource(component.GetOverrides(ctx), objectName, objectKind); found {
				return component.Name(), found
			}
		}
	}
	return "", false
}

// componentContainsResource looks through the component override list see if the resource is listed
func componentContainsResource(Overrides []installv1alpha1.Overrides, objectName string, objectKind string) bool {
	for _, override := range Overrides {
		if objectKind == constants.ConfigMapKind && override.ConfigMapRef != nil {
			if objectName == override.ConfigMapRef.Name {
				return true
			}
		}
		if objectKind == constants.SecretKind && override.SecretRef != nil {
			if objectName == override.SecretRef.Name {
				return true
			}
		}
	}
	return false
}

// UpdateVerrazzanoForInstallOverrides mutates the status subresource of Verrazzano Custom Resource specific
// to a component to cause a reconcile
func UpdateVerrazzanoForInstallOverrides(c client.Client, componentCtx spi.ComponentContext, componentName string) error {
	cr := componentCtx.ActualCR()
	// Return an error to requeue if Verrazzano Component Status hasn't been initialized
	if cr.Status.Components == nil {
		return fmt.Errorf("Components not initialized")
	}
	// Set ReconcilingGeneration to 1 to re-enter install flow
	cr.Status.Components[componentName].ReconcilingGeneration = 1
	err := c.Status().Update(context.TODO(), cr)
	if err == nil {
		return nil
	}
	return err
}

// ProcDeletedOverride checks Verrazzano CR for an override resource that has now been deleted,
// and updates the CR if the resource is found listed as an override
func ProcDeletedOverride(c client.Client, vz *installv1alpha1.Verrazzano, objectName string, objectKind string) error {

	// DefaultLogger is used since we only need to create a component context and any actual logging isn't being performed
	log := vzlog.DefaultLogger()
	ctx, err := spi.NewContext(log, c, vz, false)
	if err != nil {
		return err
	}

	compName, ok := VzContainsResource(ctx, objectName, objectKind)
	if !ok {
		return nil
	}

	if err := UpdateVerrazzanoForInstallOverrides(c, ctx, compName); err != nil {
		return err
	}
	return nil
}

// RetrieveInstallOverrideResources takes the list of Overrides and returns a list of key value pairs
//func RetrieveInstallOverrideResources(ctx spi.ComponentContext, overrides []installv1alpha1.Overrides) ([]bom.KeyValue, error) {
//	var kvs []bom.KeyValue
//	for _, override := range overrides {
//		// Check if ConfigMapRef is populated and gather helm file
//		if override.ConfigMapRef != nil {
//			// Get the ConfigMap
//			configMap := &v1.ConfigMap{}
//			selector := override.ConfigMapRef
//			nsn := types.NamespacedName{Name: selector.Name, Namespace: ctx.EffectiveCR().Namespace}
//			optional := selector.Optional
//			err := ctx.Client().Get(context.TODO(), nsn, configMap)
//			if err != nil {
//				if optional == nil || !*optional {
//					err := ctx.Log().ErrorfNewErr("Could not get Configmap %s from namespace %s: %v", nsn.Name, nsn.Namespace, err)
//					return kvs, err
//				}
//				ctx.Log().Debugf("Optional Configmap %s from namespace %s not found", nsn.Name, nsn.Namespace)
//				continue
//			}
//
//			tmpFile, err := createInstallOverrideFile(ctx, nsn, configMap.Data, selector.Key, selector.Optional)
//			if err != nil {
//				return kvs, err
//			}
//			if tmpFile != nil {
//				kvs = append(kvs, bom.KeyValue{Value: tmpFile.Name(), IsFile: true})
//			}
//		}
//		// Check if SecretRef is populated and gather helm file
//		if override.SecretRef != nil {
//			// Get the Secret
//			sec := &v1.Secret{}
//			selector := override.SecretRef
//			nsn := types.NamespacedName{Name: selector.Name, Namespace: ctx.EffectiveCR().Namespace}
//			optional := selector.Optional
//			err := ctx.Client().Get(context.TODO(), nsn, sec)
//			if err != nil {
//				if optional == nil || !*optional {
//					err := ctx.Log().ErrorfNewErr("Could not get Secret %s from namespace %s: %v", nsn.Name, nsn.Namespace, err)
//					return kvs, err
//				}
//				ctx.Log().Debugf("Optional Secret %s from namespace %s not found", nsn.Name, nsn.Namespace)
//				continue
//			}
//
//			dataStrings := map[string]string{}
//			for key, val := range sec.Data {
//				dataStrings[key] = string(val)
//			}
//			tmpFile, err := createInstallOverrideFile(ctx, nsn, dataStrings, selector.Key, selector.Optional)
//			if err != nil {
//				return kvs, err
//			}
//			if tmpFile != nil {
//				kvs = append(kvs, bom.KeyValue{Value: tmpFile.Name(), IsFile: true})
//			}
//		}
//	}
//	return kvs, nil
//}
//
//// createInstallOverrideFile takes in the data from a kubernetes resource and creates a temporary file for helm install
//func createInstallOverrideFile(ctx spi.ComponentContext, nsn types.NamespacedName, data map[string]string, dataKey string, optional *bool) (*os.File, error) {
//	var file *os.File
//
//	// Get resource data
//	fieldData, ok := data[dataKey]
//	if !ok {
//		if optional == nil || !*optional {
//			err := ctx.Log().ErrorfNewErr("Could not get Data field %s from Resource %s from namespace %s", dataKey, nsn.Name, nsn.Namespace)
//			return file, err
//		}
//		ctx.Log().Debugf("Optional Resource %s from namespace %s missing Data key %s", nsn.Name, nsn.Namespace, dataKey)
//		return file, nil
//	}
//
//	// Create the temp file for the data
//	file, err := os2.CreateTempFile(ctx.Log(), "install-overrides-*.yaml", []byte(fieldData))
//	if err != nil {
//		return file, err
//	}
//	return file, nil
//}
