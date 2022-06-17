package coherence

import (
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/helm"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/module/modules"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/secret"
	"github.com/verrazzano/verrazzano/platform-operator/internal/config"
)

const (
	ComponentName = "coherence-operator"
)

type Controller struct {
	modules.Reconciler
}

func NewComponent() modules.DelegateReconciler {
	return &Controller{
		modules.Reconciler{
			ChartDir: config.GetThirdPartyDir(),
			HelmComponent: helm.HelmComponent{
				ImagePullSecretKeyname: secret.DefaultImagePullSecretKeyName,
			},
		},
	}
}
