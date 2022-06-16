package modules

import (
	"context"
	"fmt"
	vzlogInit "github.com/verrazzano/verrazzano/pkg/log"
	vzstring "github.com/verrazzano/verrazzano/pkg/string"
	modulesv1alpha1 "github.com/verrazzano/verrazzano/platform-operator/apis/modules/v1alpha1"
	vzapi "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/helm"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	"go.uber.org/zap"
	"path/filepath"
)

const FinalizerName = "modules.finalizer.verrazzano.io"

type Reconciler struct {
	ChartDir string
	helm.HelmComponent
}

func (r *Reconciler) Reconcile(ctx spi.ComponentContext) error {
	r.InitForModule(ctx.Module())
	if ctx.Module().IsBeingDeleted() {
		return r.Uninstall(ctx)
	}
	if err := r.PreUpgrade(ctx); err != nil {
		return err
	}
	if err := r.Install(ctx); err != nil {
		return err
	}
	if err := r.PostUpgrade(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Reconciler) InitForModule(module *modulesv1alpha1.Module) {
	chart := module.Spec.Installer.HelmChart
	if chart != nil {
		r.ReleaseName = chart.Name
		r.JSONName = chart.Name
		r.HelmComponent.ChartDir = filepath.Join(r.ChartDir, chart.Repository.Path)
		r.ChartNamespace = chart.Namespace
		r.IgnoreNamespaceOverride = true
		r.GetInstallOverridesFunc = func(_ *vzapi.Verrazzano) []vzapi.Overrides {
			return chart.InstallOverrides.ValueOverrides
		}
	}
}

func (r *Reconciler) PreUpgrade(ctx spi.ComponentContext) error {
	return addFinalizer(ctx)
}

func (r *Reconciler) Ready(_ spi.ComponentContext) bool {
	return true
}

func (r *Reconciler) PostUpgrade(_ spi.ComponentContext) error {
	return nil
}

func (r *Reconciler) UpdatePhase(ctx spi.ComponentContext, status modulesv1alpha1.ModulePhase) error {
	return ctx.Client().Update(context.TODO(), ctx.Module())
}

//Uninstall cleans up the Helm Chart and removes the Module finalizer so Kubernetes can clean the resource
func (r *Reconciler) Uninstall(ctx spi.ComponentContext) error {
	if err := r.HelmComponent.Uninstall(ctx); err != nil {
		return err
	}
	return removeFinalizer(ctx)
}

func removeFinalizer(ctx spi.ComponentContext) error {
	module := ctx.Module()
	if needsFinalizerRemoval(module) {
		module.Finalizers = vzstring.RemoveStringFromSlice(module.Finalizers, FinalizerName)
		err := ctx.Client().Update(context.TODO(), module)
		return vzlogInit.ConflictWithLog(fmt.Sprintf("Failed to remove finalizer from module %s/%s", module.Namespace, module.Name), err, zap.S())
	}
	return nil
}

func addFinalizer(ctx spi.ComponentContext) error {
	module := ctx.Module()
	if needsFinalizer(module) {
		module.Finalizers = append(module.Finalizers, FinalizerName)
		err := ctx.Client().Update(context.TODO(), module)
		_, err = vzlogInit.IgnoreConflictWithLog(fmt.Sprintf("Failed to add finalizer to ingress trait %s", module.Name), err, zap.S())
		return err
	}
	return nil
}

func needsFinalizer(module *modulesv1alpha1.Module) bool {
	return module.GetDeletionTimestamp().IsZero() && !vzstring.SliceContainsString(module.Finalizers, FinalizerName)
}

func needsFinalizerRemoval(module *modulesv1alpha1.Module) bool {
	return !needsFinalizer(module)
}
