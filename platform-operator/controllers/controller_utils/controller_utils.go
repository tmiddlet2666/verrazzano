package controller_utils

import (
	"context"
	"fmt"
	spi2 "github.com/verrazzano/verrazzano/pkg/controller/errors"
	"github.com/verrazzano/verrazzano/pkg/log/vzlog"
	"github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/controllers"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/vzinstance"
	"k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func UpdateComponentStatus(c client.Client, compContext spi.ComponentContext, message string, conditionType v1alpha1.ConditionType) error {
	t := time.Now().UTC()
	condition := v1alpha1.Condition{
		Type:    conditionType,
		Status:  v1.ConditionTrue,
		Message: message,
		LastTransitionTime: fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second()),
	}

	componentName := compContext.GetComponent()
	cr := compContext.ActualCR()
	log := compContext.Log()

	if cr.Status.Components == nil {
		cr.Status.Components = make(map[string]*v1alpha1.ComponentStatusDetails)
	}
	componentStatus := cr.Status.Components[componentName]
	if componentStatus == nil {
		componentStatus = &v1alpha1.ComponentStatusDetails{
			Name: componentName,
		}
		cr.Status.Components[componentName] = componentStatus
	}
	if conditionType == v1alpha1.CondInstallComplete {
		cr.Status.VerrazzanoInstance = vzinstance.GetInstanceInfo(compContext)
		if componentStatus.ReconcilingGeneration > 0 {
			componentStatus.LastReconciledGeneration = componentStatus.ReconcilingGeneration
			componentStatus.ReconcilingGeneration = 0
		} else {
			componentStatus.LastReconciledGeneration = cr.Generation
		}
	} else {
		if componentStatus.ReconcilingGeneration == 0 {
			componentStatus.ReconcilingGeneration = cr.Generation
		}
	}
	componentStatus.Conditions = verrazzano.AppendConditionIfNecessary(log, componentStatus, condition)

	// Set the state of resource
	componentStatus.State = controllers.CheckCondtitionType(conditionType)

	// Update the status
	return UpdateVerrazzanoStatus(compContext.Client(), log, cr)
}

func UpdateVerrazzanoStatus(c client.Client, log vzlog.VerrazzanoLogger, vz *v1alpha1.Verrazzano) error {
	err := c.Status().Update(context.TODO(), vz)
	if err == nil {
		return nil
	}
	if spi2.IsUpdateConflict(err) {
		log.Debugf("Requeuing to get a fresh copy of the Verrazzano resource since the current one is outdated.")
	} else {
		log.Errorf("Failed to update Verrazzano resource :v", err)
	}
	// Return error so that reconcile gets called again
	return err
}
