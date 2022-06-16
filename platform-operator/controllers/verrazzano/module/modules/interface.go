package modules

import "github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"

const ControllerLabel = "verrazzano.io/module"

type DelegateReconciler interface {
	Reconcile(ctx spi.ComponentContext) error
}
