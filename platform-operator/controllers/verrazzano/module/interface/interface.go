package _interface

import (
	"context"
	"github.com/verrazzano/verrazzano/platform-operator/apis/modules/v1alpha1"
	clipkg "sigs.k8s.io/controller-runtime/pkg/client"
)

const ControllerLabel = "verrazzano.io/module"

type DelegateReconciler interface {
	Reconcile(ctx context.Context, client clipkg.Client, module *v1alpha1.Module) error
}

type BaseReconciler struct{}

func (b BaseReconciler) Reconcile(ctx context.Context, client clipkg.Client, module *v1alpha1.Module) error {
	if err := b.PreHook(ctx, client, module); err != nil {
		return err
	}
	if err := b.Install(ctx, client, module); err != nil {
		return err
	}
	return b.PostHook(ctx, client, module)
}

func (b BaseReconciler) PreHook(ctx context.Context, client clipkg.Client, module *v1alpha1.Module) error {
	return nil
}

func (b BaseReconciler) PostHook(ctx context.Context, client clipkg.Client, module *v1alpha1.Module) error {
	return nil
}

func (b BaseReconciler) Install(ctx context.Context, client clipkg.Client, module *v1alpha1.Module) error {
	return nil
}
