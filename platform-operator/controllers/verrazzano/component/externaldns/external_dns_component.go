// Copyright (c) 2021, 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package externaldns

import (
	"fmt"
	modulesv1alpha1 "github.com/verrazzano/verrazzano/platform-operator/apis/modules/v1alpha1"
	vzapi "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/module/modules"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/module/reconciler"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/common"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/helm"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	"github.com/verrazzano/verrazzano/platform-operator/internal/config"
)

// ComponentName is the name of the component
const ComponentName = "external-dns"

// ComponentNamespace is the namespace of the component
const ComponentNamespace = "cert-manager"

const overrideFile = "external-dns-values.yaml"
const ConfigMapName = "external-dns-vz-config"

type externalDNSComponent struct {
	helm.HelmComponent
}

// Verify that nginxComponent implements Component
var _ spi.Component = externalDNSComponent{}

func NewComponent(module *modulesv1alpha1.Module) modules.DelegateReconciler {
	h := helm.HelmComponent{
		ChartDir:               config.GetThirdPartyDir(),
		ImagePullSecretKeyname: imagePullSecretHelmKey,
		AppendOverridesFunc:    AppendOverrides,
	}
	helm.SetForModule(&h, module)

	return &reconciler.Reconciler{
		ModuleComponent: externalDNSComponent{
			h,
		},
	}
}

func (e externalDNSComponent) PreUpgrade(compContext spi.ComponentContext) error {
	return common.ApplyOverride(compContext, overrideFile)
}

func (e externalDNSComponent) PreInstall(compContext spi.ComponentContext) error {
	return preInstall(compContext)
}

func (e externalDNSComponent) IsOperatorInstallSupported() bool {
	return false
}

func (e externalDNSComponent) Name() string {
	if e.HelmComponent.ReleaseName == "" {
		return ComponentName
	}
	return e.HelmComponent.ReleaseName
}

func (e externalDNSComponent) IsReady(ctx spi.ComponentContext) bool {
	if e.HelmComponent.IsReady(ctx) {
		return isExternalDNSReady(ctx)
	}
	return false
}

func (e externalDNSComponent) IsEnabled(effectiveCR *vzapi.Verrazzano) bool {
	dns := effectiveCR.Spec.Components.DNS
	if dns != nil && dns.OCI != nil {
		return true
	}
	return false
}

// ValidateUpdate checks if the specified new Verrazzano CR is valid for this component to be updated
func (e externalDNSComponent) ValidateUpdate(old *vzapi.Verrazzano, new *vzapi.Verrazzano) error {
	// Do not allow any changes except to enable the component post-install
	if e.IsEnabled(old) && !e.IsEnabled(new) {
		return fmt.Errorf("Disabling an existing OCI DNS configuration is not allowed")
	}
	return e.HelmComponent.ValidateUpdate(old, new)
}

// MonitorOverrides checks whether monitoring of install overrides is enabled or not
func (e externalDNSComponent) MonitorOverrides(ctx spi.ComponentContext) bool {
	if ctx.EffectiveCR().Spec.Components.DNS != nil {
		if ctx.EffectiveCR().Spec.Components.DNS.MonitorChanges != nil {
			return *ctx.EffectiveCR().Spec.Components.DNS.MonitorChanges
		}
		return true
	}
	return false
}
