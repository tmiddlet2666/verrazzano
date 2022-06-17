package modules

import (
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	clipkg "sigs.k8s.io/controller-runtime/pkg/client"
)

const ControllerLabel = "verrazzano.io/module"

type DelegateReconciler interface {
	Reconcile(ctx spi.ComponentContext) error
	SetStatusWriter(statusWriter clipkg.StatusWriter)
}
