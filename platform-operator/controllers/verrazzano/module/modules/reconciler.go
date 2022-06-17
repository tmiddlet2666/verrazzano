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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const FinalizerName = "modules.finalizer.verrazzano.io"

type Reconciler struct {
	client.StatusWriter
	ChartDir string
	helm.HelmComponent
}

func (r *Reconciler) SetStatusWriter(writer client.StatusWriter) {
	r.StatusWriter = writer
}

func (r *Reconciler) doReconcile(ctx spi.ComponentContext) error {
	if ctx.Module().Status.Phase == nil {
		ctx.Module().SetPhase(modulesv1alpha1.PhasePending)
	}
	phase := *ctx.Module().Status.Phase
	switch phase {
	case modulesv1alpha1.PhasePending:
		return r.PendingPhase(ctx)
	case modulesv1alpha1.PhaseInstalling:
		return r.InstallingPhase(ctx)
	case modulesv1alpha1.PhaseReconciling:
		return r.ReconcilingPhase(ctx)
	case modulesv1alpha1.PhaseNotReady:
		return r.NotReadyPhase(ctx)
	case modulesv1alpha1.PhaseReady:
		return r.ReadyPhase(ctx)
	}
	return nil
}

func (r *Reconciler) Reconcile(ctx spi.ComponentContext) error {
	r.InitForModule(ctx.Module())
	// Delete module if it is being deleted
	if ctx.Module().IsBeingDeleted() {
		if err := r.UpdatePhaseOfModule(ctx, modulesv1alpha1.PhaseDeleting); err != nil {
			return err
		}
		return r.Uninstall(ctx)
	}
	return r.doReconcile(ctx)
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

func (r *Reconciler) PendingPhase(ctx spi.ComponentContext) error {
	if err := addFinalizer(ctx); err != nil {
		return err
	}
	return r.UpdatePhaseOfModule(ctx, modulesv1alpha1.PhaseInstalling)
}

func (r *Reconciler) InstallingPhase(ctx spi.ComponentContext) error {
	if err := r.HelmComponent.Install(ctx); err != nil {
		return err
	}
	return r.UpdatePhaseOfModule(ctx, modulesv1alpha1.PhaseReconciling)
}

func (r *Reconciler) IsReady(_ spi.ComponentContext) bool {
	return true
}

func (r *Reconciler) NotReadyPhase(ctx spi.ComponentContext) error {
	if r.IsReady(ctx) {
		module := ctx.Module()
		module.Status.ObservedGeneration = module.Generation
		return r.UpdatePhaseOfModule(ctx, modulesv1alpha1.PhaseReady)
	}
	return nil
}

//ReadyPhase reconciles put the Module back to pending state if the generation has changed
func (r *Reconciler) ReadyPhase(ctx spi.ComponentContext) error {
	module := ctx.Module()
	if module.Status.ObservedGeneration != module.Generation {
		return r.UpdatePhaseOfModule(ctx, modulesv1alpha1.PhasePending)
	}
	return nil
}

func (r *Reconciler) ReconcilingPhase(ctx spi.ComponentContext) error {
	return r.UpdatePhaseOfModule(ctx, modulesv1alpha1.PhaseNotReady)
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

func (r *Reconciler) UpdatePhaseOfModule(ctx spi.ComponentContext, phase modulesv1alpha1.ModulePhase) error {
	ctx.Module().SetPhase(phase)
	return r.StatusWriter.Update(context.TODO(), ctx.Module())
}
