// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package reconciler

import (
	"context"
	"fmt"
	modulesv1alpha1 "github.com/verrazzano/verrazzano/platform-operator/apis/modules/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	ctrlUtils "github.com/verrazzano/verrazzano/platform-operator/controllers/controller_utils"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"time"
)

//UpdateStatus configures the Module's status based on the passed in state and then updates the Module on the cluster
func (r *Reconciler) UpdateStatus(ctx spi.ComponentContext, moduleCondition modulesv1alpha1.ModuleCondition) error {
	phase := modulesv1alpha1.Phase(moduleCondition)
	// Update the Module's Phase
	ctx.Module().SetPhase(phase)
	// Append a new moduleCondition, if applicable
	appendCondition(ctx.Module(), string(phase), moduleCondition)

	// update the Verrazzano CR component status to align with the module status
	if err := ctrlUtils.UpdateComponentStatus(ctx.Client(), ctx, string(phase), convertModuleConditiontoCondition(moduleCondition)); err != nil {
		return err
	}

	return r.doStatusUpdate(ctx)
}

func NeedsReconcile(ctx spi.ComponentContext) bool {
	return ctx.Module().Status.ObservedGeneration != ctx.Module().Generation
}

func NewCondition(message string, condition modulesv1alpha1.ModuleCondition) modulesv1alpha1.Condition {
	t := time.Now().UTC()
	return modulesv1alpha1.Condition{
		Type:    condition,
		Message: message,
		Status:  corev1.ConditionTrue,
		LastTransitionTime: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second()),
	}
}

func (r *Reconciler) doStatusUpdate(ctx spi.ComponentContext) error {
	module := ctx.Module()
	err := r.StatusWriter.Update(context.TODO(), module)
	if err == nil {
		return err
	}
	if k8serrors.IsConflict(err) {
		ctx.Log().Debugf("Update conflict for Module %s: %v", module.Name, err)
	} else {
		ctx.Log().Errorf("Failed to update Module %s :v", module.Name, err)
	}
	// Return error so that reconcile gets called again
	return err
}

func appendCondition(module *modulesv1alpha1.Module, message string, condition modulesv1alpha1.ModuleCondition) {
	conditions := module.Status.Conditions
	lastCondition := conditions[len(conditions)-1]
	newCondition := NewCondition(message, condition)
	// Only update the conditions if there is a notable change between the last update
	if needsConditionUpdate(lastCondition, newCondition) {
		// Delete oldest condition if at tracking limit
		if len(conditions) > modulesv1alpha1.ConditionArrayLimit {
			conditions = conditions[1:]
		}
		module.Status.Conditions = append(conditions, newCondition)
	}
}

//needsConditionUpdate checks if the condition needs an update
func needsConditionUpdate(last, new modulesv1alpha1.Condition) bool {
	return last.Type != new.Type && last.Message != new.Message
}

// convertModuleConditiontoCondition converts ModuleCondition types to ConditionType types
// this will then get converted to CompStateType in updateComponentStatus by the CheckCondtitionType function
// return nil if ModuleCondition is unknown
func convertModuleConditiontoCondition(moduleCondion modulesv1alpha1.ModuleCondition) v1alpha1.ConditionType {
	switch moduleCondion {
	// install ConditionTypes
	case modulesv1alpha1.CondPreInstall:
		return v1alpha1.CondPreInstall
	case modulesv1alpha1.CondInstallStarted:
		return v1alpha1.CondInstallStarted
	case modulesv1alpha1.CondInstallComplete:
		return v1alpha1.CondInstallComplete
	// upgrade ConditionTypes
	case modulesv1alpha1.CondPreUpgrade, modulesv1alpha1.CondUpgradeStarted:
		return v1alpha1.CondUpgradeStarted
	case modulesv1alpha1.CondUpgradeComplete:
		return v1alpha1.CondUpgradeComplete
	}
	// otherwise return UninstallStarted
	return v1alpha1.CondUninstallStarted
}
