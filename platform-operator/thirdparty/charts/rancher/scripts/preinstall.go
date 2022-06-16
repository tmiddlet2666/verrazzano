// Copyright (c) 2021, 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//const (
//	// note: VZ-5241 In Rancher 2.6.3 the agent was moved from cattle-fleet-system ns
//	// to a new cattle-fleet-local-system ns, the rancher-operator-system ns was
//	// removed, and the rancher-operator is no longer deployed
//	FleetSystemNamespace      = "cattle-fleet-system"
//	FleetLocalSystemNamespace = "cattle-fleet-local-system"
//	defaultSecretNamespace    = "cert-manager"
//	namespaceLabelKey         = "verrazzano.io/namespace"
//	rancherTLSSecretName      = "tls-ca"
//	defaultVerrazzanoName     = "verrazzano-ca-certificate-secret"
//	fleetAgentDeployment      = "fleet-agent"
//	fleetControllerDeployment = "fleet-controller"
//	gitjobDeployment          = "gitjob"
//	rancherWebhookDeployment  = "rancher-webhook"
//	letsEncryptTLSSource       = "letsEncrypt"
//	caTLSSource                = "secret"
//	caCertsPem                 = "cacerts.pem"
//	caCert                     = "ca.crt"
//	privateCAValue             = "true"
//	useBundledSystemChartValue = "true"
//)

func main() {
	zaplog := zap.S()
	fmt.Println("Entered main print")
	log.Println("Entered main log")
	zaplog.Info("Entered main zap")
	config, err := ctrl.GetConfig()
	if err != nil {
		zaplog.Errorf("Failed to get kubeconfig: %v", err)
		os.Exit(1)
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		zaplog.Errorf("Failed to create clientset: %s", err.Error())
		os.Exit(1)
	}
	if err := createCattleSystemNamespace(zaplog, c); err != nil {
		zaplog.Errorf("Failed creating cattle-system namespace: %s", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
	//if err := copyDefaultCACertificate(zaplog, c, vz); err != nil {
	//	zaplog.Errorf("Failed copying default CA certificate: %s", err.Error())
	//	os.Exit(1)
	//}
}

// createCattleSystemNamespace creates the cattle-system namespace if it does not exist
func createCattleSystemNamespace(_ *zap.SugaredLogger, c *kubernetes.Clientset) error {
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   common.CattleSystem,
			Labels: map[string]string{"test2": "val"},
		},
	}
	fmt.Printf("Creating %s namespace", common.CattleSystem)
	_, err := c.CoreV1().Namespaces().Update(context.TODO(), namespace, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

////copyDefaultCACertificate copies the defaultVerrazzanoName TLS Secret to the ComponentNamespace for use by Rancher
////This method will only copy defaultVerrazzanoName if default CA certificates are being used.
//func copyDefaultCACertificate(log vzlog.VerrazzanoLogger, c client.Client, vz *vzapi.Verrazzano) error {
//	cm := vz.Spec.Components.CertManager
//	if isUsingDefaultCACertificate(cm) {
//		namespacedName := types.NamespacedName{Namespace: defaultSecretNamespace, Name: defaultVerrazzanoName}
//		defaultSecret := &v1.Secret{}
//		if err := c.Get(context.TODO(), namespacedName, defaultSecret); err != nil {
//			return err
//		}
//		if len(defaultSecret.Data[caCert]) < 1 {
//			return log.ErrorfNewErr("Failed, secret %s/%s does not have a value for %s", defaultSecretNamespace, defaultVerrazzanoName, caCert)
//		}
//		rancherCaSecret := &v1.Secret{
//			ObjectMeta: metav1.ObjectMeta{
//				Namespace: common.CattleSystem,
//				Name:      rancherTLSSecretName,
//			},
//		}
//		log.Debugf("Copying default Verrazzano secret to Rancher namespace")
//		if _, err := controllerruntime.CreateOrUpdate(context.TODO(), c, rancherCaSecret, func() error {
//			rancherCaSecret.Data = map[string][]byte{
//				caCertsPem: defaultSecret.Data[caCert],
//			}
//			return nil
//		}); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func isUsingDefaultCACertificate(cm *vzapi.CertManagerComponent) bool {
//	return cm != nil &&
//		cm.Certificate.CA != vzapi.CA{} &&
//		cm.Certificate.CA.SecretName == defaultVerrazzanoName &&
//		cm.Certificate.CA.ClusterResourceNamespace == defaultSecretNamespace
//}
