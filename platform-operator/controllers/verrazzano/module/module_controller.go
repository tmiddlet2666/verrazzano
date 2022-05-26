package module

import (
	"context"
	"github.com/verrazzano/verrazzano/platform-operator/apis/modules/v1alpha1"
	clipkg "sigs.k8s.io/controller-runtime/pkg/client"
)

type ModuleReconciler struct {
	Controllers []ModuleController
}

type ModuleController interface {
	Reconcile(ctx context.Context, client clipkg.Client, module *v1alpha1.Module) error
}

type BaseModuleController struct{}

func (b BaseModuleController) Reconcile(ctx context.Context, client clipkg.Client, module *v1alpha1.Module) error {
	return nil
}
