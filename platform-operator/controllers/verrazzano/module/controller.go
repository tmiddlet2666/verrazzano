package module

import (
	"context"
	"github.com/verrazzano/verrazzano/application-operator/controllers/clusters"
	"github.com/verrazzano/verrazzano/pkg/log/vzlog"
	modulesv1alpha1 "github.com/verrazzano/verrazzano/platform-operator/apis/modules/v1alpha1"
	vzapi "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/module/modules"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/module/modules/coherence"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/module/modules/rancher"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

var delegates = map[string]func() modules.DelegateReconciler{
	coherence.ComponentName: coherence.NewComponent,
	rancher.ComponentName:   rancher.NewComponent,
}

type Reconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Controller controller.Controller
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&modulesv1alpha1.Module{}).
		Complete(r)
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	verrazzanos := &vzapi.VerrazzanoList{}
	if err := r.List(ctx, verrazzanos); err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		zap.S().Errorf("Failed to get Verrazzanos %s/%s", req.Namespace, req.Name)
		return clusters.NewRequeueWithDelay(), err
	}

	// Get the module for the request
	module := &modulesv1alpha1.Module{}
	if err := r.Get(ctx, req.NamespacedName, module); err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		zap.S().Errorf("Failed to get Module %s/%s", req.Namespace, req.Name)
		return clusters.NewRequeueWithDelay(), err
	}
	// Get the resource logger needed to log message using 'progress' and 'once' methods
	log, err := vzlog.EnsureResourceLogger(&vzlog.ResourceConfig{
		Name:           module.Name,
		Namespace:      module.Namespace,
		ID:             string(module.UID),
		Generation:     module.Generation,
		ControllerName: "verrazzano",
	})
	if err != nil {
		zap.S().Errorf("Failed to create controller logger for Module controller: %v", err)
		return clusters.NewRequeueWithDelay(), err
	}

	if module.Generation == module.Status.ObservedGeneration {
		log.Debugf("Skipping module %s reconcile, observed generation has not change", module.Name)
		return ctrl.Result{}, nil
	}

	// Unknown module controller cannot be handled
	delegate := getDelegateController(module)
	if delegate == nil {
		return ctrl.Result{}, nil
	}
	moduleCtx, err := spi.NewModuleContext(log, r.Client, &verrazzanos.Items[0], module, false)
	if err != nil {
		log.Errorf("Failed to create module context: %v", err)
		return clusters.NewRequeueWithDelay(), err
	}
	delegate.SetStatusWriter(r.Status())
	if err := delegate.Reconcile(moduleCtx); err != nil {
		log.Errorf("Failed to reconcile module %s/%s: %v", module.Name, module.Namespace, err)
		return clusters.NewRequeueWithDelay(), err
	}
	return ctrl.Result{}, nil
}

func getDelegateController(module *modulesv1alpha1.Module) modules.DelegateReconciler {
	newDelegate := delegates[module.ObjectMeta.Labels[modules.ControllerLabel]]
	if newDelegate == nil {
		return nil
	}
	return newDelegate()
}
