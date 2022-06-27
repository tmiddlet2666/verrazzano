// Copyright (c) 2021, 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package weblogic

import (
	"fmt"
	modulesv1alpha1 "github.com/verrazzano/verrazzano/platform-operator/apis/modules/v1alpha1"
	vzapi "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/module/modules"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/module/reconciler"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/common"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/secret"
	"github.com/verrazzano/verrazzano/platform-operator/internal/config"
	"github.com/verrazzano/verrazzano/platform-operator/internal/vzconfig"

	"github.com/verrazzano/verrazzano/platform-operator/constants"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/helm"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
)

// ComponentName is the name of the component
const ComponentName = "weblogic-operator"

// ComponentNamespace is the namespace of the component
const ComponentNamespace = constants.VerrazzanoSystemNamespace

// ComponentJSONName is the josn name of the verrazzano component in CRD
const ComponentJSONName = "weblogicOperator"

const ConfigMapName = "weblogic-operator-vz-config"

const overridesFile = "weblogic-values.yaml"

type weblogicComponent struct {
	helm.HelmComponent
}

func NewComponent(module *modulesv1alpha1.Module) modules.DelegateReconciler {
	h := helm.HelmComponent{
		ChartDir:               config.GetThirdPartyDir(),
		ImagePullSecretKeyname: secret.DefaultImagePullSecretKeyName,
		PreInstallFunc:         WeblogicOperatorPreInstall,
		AppendOverridesFunc:    AppendWeblogicOperatorOverrides,
	}
	helm.SetForModule(&h, module)

	return &reconciler.Reconciler{
		ModuleComponent: weblogicComponent{
			h,
		},
	}
}

func (c weblogicComponent) PreInstall(ctx spi.ComponentContext) error {
	return common.ApplyOverride(ctx, overridesFile)
}

func (c weblogicComponent) PreUpgrade(ctx spi.ComponentContext) error {
	return common.ApplyOverride(ctx, overridesFile)
}

// IsEnabled WebLogic-specific enabled check for installation
func (c weblogicComponent) IsEnabled(effectiveCR *vzapi.Verrazzano) bool {
	return vzconfig.IsWeblogicOperatorEnabled(effectiveCR)
}

// ValidateUpdate checks if the specified new Verrazzano CR is valid for this component to be updated
func (c weblogicComponent) ValidateUpdate(old *vzapi.Verrazzano, new *vzapi.Verrazzano) error {
	// Do not allow disabling of any component post-install for now
	if c.IsEnabled(old) && !c.IsEnabled(new) {
		return fmt.Errorf("Disabling component %s is not allowed", ComponentJSONName)
	}
	return c.HelmComponent.ValidateUpdate(old, new)
}

// IsReady component check
func (c weblogicComponent) IsReady(ctx spi.ComponentContext) bool {
	if c.HelmComponent.IsReady(ctx) {
		return isWeblogicOperatorReady(ctx)
	}
	return false
}

// MonitorOverrides checks whether monitoring of install overrides is enabled or not
func (c weblogicComponent) MonitorOverrides(ctx spi.ComponentContext) bool {
	if ctx.EffectiveCR().Spec.Components.WebLogicOperator != nil {
		if ctx.EffectiveCR().Spec.Components.WebLogicOperator.MonitorChanges != nil {
			return *ctx.EffectiveCR().Spec.Components.WebLogicOperator.MonitorChanges
		}
		return true
	}
	return false
}

func (c weblogicComponent) IsOperatorInstallSupported() bool {
	return false
}

func (c weblogicComponent) Name() string {
	if c.HelmComponent.ReleaseName == "" {
		return ComponentName
	}
	return c.HelmComponent.ReleaseName
}
