package certmanager

import (
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/helm"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/module/modules"
	"github.com/verrazzano/verrazzano/platform-operator/internal/config"
)

const (
	ComponentName = "cert-manager"
)

type Controller struct {
	modules.Reconciler
}

func NewComponent() *Controller {
	return &Controller{
		modules.Reconciler{
			ChartDir: config.GetThirdPartyDir(),
			HelmComponent: helm.HelmComponent{
				ImagePullSecretKeyname: "global.imagePullSecrets[0].name",
				AppendOverridesFunc:    AppendOverrides,
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
