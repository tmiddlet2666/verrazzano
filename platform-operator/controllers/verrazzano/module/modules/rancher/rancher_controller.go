package rancher

import (
	vzapi "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/helm"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/module/modules"
)

const (
	ComponentName = "rancher"
)

type Controller struct {
	modules.Reconciler
}

func NewComponent() *Controller {
	return &Controller{
		modules.Reconciler{
			HelmComponent: helm.HelmComponent{
				GetInstallOverridesFunc: func(vz *vzapi.Verrazzano) []vzapi.Overrides {
					return nil
				},
			},
		},
	}
}

func (r *Controller) PreHook(ctx spi.ComponentContext) error {
	return nil
}

func (r *Controller) PostHook(ctx spi.ComponentContext) error {
	return nil
}
