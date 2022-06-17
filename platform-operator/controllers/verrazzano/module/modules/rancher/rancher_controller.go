package rancher

import (
	"github.com/verrazzano/verrazzano/platform-operator/constants"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/common"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/helm"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/module/modules"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/secret"
	"github.com/verrazzano/verrazzano/platform-operator/internal/config"
	"k8s.io/apimachinery/pkg/types"
)

const (
	ComponentName      = "rancher"
	ComponentNamespace = "cattle-system"
)

var certificates = []types.NamespacedName{
	{Name: "tls-rancher-ingress", Namespace: ComponentNamespace},
}

type Controller struct {
	modules.Reconciler
}

func NewComponent() modules.DelegateReconciler {
	return &Controller{
		modules.Reconciler{
			ChartDir: config.GetThirdPartyDir(),
			HelmComponent: helm.HelmComponent{
				ImagePullSecretKeyname: secret.DefaultImagePullSecretKeyName,
				AppendOverridesFunc:    AppendOverrides,
				Certificates:           certificates,
				IngressNames: []types.NamespacedName{
					{
						Namespace: ComponentNamespace,
						Name:      constants.RancherIngress,
					},
				},
			},
		},
	}
}

func (c *Controller) IsReady(ctx spi.ComponentContext) bool {
	if c.HelmComponent.IsReady(ctx) {
		return isRancherReady(ctx)
	}
	return false
}

func (c *Controller) PendingPhase(ctx spi.ComponentContext) error {
	vz := ctx.EffectiveCR()
	client := ctx.Client()
	log := ctx.Log()
	if err := createCattleSystemNamespace(log, client); err != nil {
		log.ErrorfThrottledNewErr("Failed creating cattle-system namespace: %s", err.Error())
		return err
	}
	if err := copyDefaultCACertificate(log, client, vz); err != nil {
		log.ErrorfThrottledNewErr("Failed copying default CA certificate: %s", err.Error())
		return err
	}
	return nil
}

func (c *Controller) ReconcilingPhase(ctx spi.ComponentContext) error {
	client := ctx.Client()
	log := ctx.Log()
	vz := ctx.EffectiveCR()
	// Set MKNOD Cap on Rancher deployment
	if err := patchRancherDeployment(client); err != nil {
		return log.ErrorfThrottledNewErr("Failed patching Rancher deployment: %s", err.Error())
	}
	log.Debugf("Patched Rancher deployment to support MKNOD")
	// Annotate Rancher ingress for NGINX/TLS
	if err := patchRancherIngress(client, ctx.EffectiveCR()); err != nil {
		return log.ErrorfThrottledNewErr("Failed patching Rancher ingress: %s", err.Error())
	}
	log.Debugf("Patched Rancher ingress")

	if err := createAdminSecretIfNotExists(log, client); err != nil {
		return log.ErrorfThrottledNewErr("Failed creating Rancher admin secret: %s", err.Error())
	}
	password, err := common.GetAdminSecret(client)
	if err != nil {
		return log.ErrorfThrottledNewErr("Failed getting Rancher admin secret: %s", err.Error())
	}
	rancherHostName, err := getRancherHostname(client, vz)
	if err != nil {
		return log.ErrorfThrottledNewErr("Failed getting Rancher hostname: %s", err.Error())
	}

	rest, err := common.NewClient(client, rancherHostName, password)
	if err != nil {
		return log.ErrorfThrottledNewErr("Failed getting Rancher client: %s", err.Error())
	}
	if err := rest.SetAccessToken(); err != nil {
		return log.ErrorfThrottledNewErr("Failed setting Rancher access token: %s", err.Error())
	}
	if err := rest.PutServerURL(); err != nil {
		return log.ErrorfThrottledNewErr("Failed setting Rancher server URL: %s", err.Error())
	}
	if err := removeBootstrapSecretIfExists(log, client, ctx.Module().ChartNamespace()); err != nil {
		return log.ErrorfThrottledNewErr("Failed removing Rancher bootstrap secret: %s", err.Error())
	}
	if err := c.HelmComponent.PostInstall(ctx); err != nil {
		return log.ErrorfThrottledNewErr("Failed helm component post install: %s", err.Error())
	}
	return nil
}
