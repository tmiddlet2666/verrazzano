package rancher

import (
	"context"
	"github.com/verrazzano/verrazzano/pkg/log/vzlog"
	vzapi "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	clipkg "sigs.k8s.io/controller-runtime/pkg/client"
)

// createCattleSystemNamespace creates the cattle-system namespace if it does not exist
func createCattleSystemNamespace(log vzlog.VerrazzanoLogger, client clipkg.Client) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: common.CattleSystem,
		},
	}
	log.Debugf("Creating %s namespace", common.CattleSystem)
	if _, err := controllerruntime.CreateOrUpdate(context.TODO(), client, namespace, func() error {
		log.Debugf("Ensuring %s label is present on %s namespace", namespaceLabelKey, common.CattleSystem)
		if namespace.Labels == nil {
			namespace.Labels = map[string]string{}
		}
		namespace.Labels[namespaceLabelKey] = common.RancherName
		return nil
	}); err != nil {
		return err
	}

	return nil
}

//copyDefaultCACertificate copies the defaultVerrazzanoName TLS Secret to the ComponentNamespace for use by Rancher
//This method will only copy defaultVerrazzanoName if default CA certificates are being used.
func copyDefaultCACertificate(log vzlog.VerrazzanoLogger, client clipkg.Client, vz *vzapi.Verrazzano) error {
	cm := vz.Spec.Components.CertManager
	if isUsingDefaultCACertificate(cm) {
		namespacedName := types.NamespacedName{Namespace: defaultSecretNamespace, Name: defaultVerrazzanoName}
		defaultSecret := &corev1.Secret{}
		if err := client.Get(context.TODO(), namespacedName, defaultSecret); err != nil {
			return err
		}
		if len(defaultSecret.Data[caCert]) < 1 {
			return log.ErrorfNewErr("Failed, secret %s/%s does not have a value for %s", defaultSecretNamespace, defaultVerrazzanoName, caCert)
		}
		rancherCaSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: common.CattleSystem,
				Name:      rancherTLSSecretName,
			},
		}
		log.Debugf("Copying default Verrazzano secret to Rancher namespace")
		if _, err := controllerruntime.CreateOrUpdate(context.TODO(), client, rancherCaSecret, func() error {
			rancherCaSecret.Data = map[string][]byte{
				caCertsPem: defaultSecret.Data[caCert],
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func isUsingDefaultCACertificate(cm *vzapi.CertManagerComponent) bool {
	return cm != nil &&
		cm.Certificate.CA != vzapi.CA{} &&
		cm.Certificate.CA.SecretName == defaultVerrazzanoName &&
		cm.Certificate.CA.ClusterResourceNamespace == defaultSecretNamespace
}
