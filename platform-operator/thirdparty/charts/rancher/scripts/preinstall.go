// Copyright (c) 2021, 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package main

import (
	"context"
	"errors"
	"fmt"
	vzapi "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/constants"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// note: VZ-5241 In Rancher 2.6.3 the agent was moved from cattle-fleet-system ns
	// to a new cattle-fleet-local-system ns, the rancher-operator-system ns was
	// removed, and the rancher-operator is no longer deployed
	defaultSecretNamespace = "cert-manager"
	rancherTLSSecretName   = "tls-ca"
	defaultVerrazzanoName  = "verrazzano-ca-certificate-secret"
	caCertsPem             = "cacerts.pem"
	caCert                 = "ca.crt"
)

func main() {
	log.Println("Entered main log")
	config, err := ctrl.GetConfig()
	if err != nil {
		log.Printf("Failed to get kubeconfig: %s", err.Error())
		os.Exit(1)
	}
	c, err := client.New(config, client.Options{})
	if err != nil {
		log.Printf("Failed to create client: %s", err.Error())
		os.Exit(1)
	}
	vz := &vzapi.Verrazzano{}
	err = c.Get(context.TODO(), client.ObjectKey{Namespace: constants.DefaultNamespace}, vz)
	if err != nil {
		log.Printf("Failed to get Verrazzano: %s", err.Error())
		os.Exit(1)
	}
	if err := copyDefaultCACertificate(c, vz); err != nil {
		log.Printf("Failed copying default CA certificate: %s", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

//copyDefaultCACertificate copies the defaultVerrazzanoName TLS Secret to the ComponentNamespace for use by Rancher
//This method will only copy defaultVerrazzanoName if default CA certificates are being used.
func copyDefaultCACertificate(c client.Client, vz *vzapi.Verrazzano) error {
	cm := vz.Spec.Components.CertManager
	if isUsingDefaultCACertificate(cm) {
		namespacedName := types.NamespacedName{Namespace: defaultSecretNamespace, Name: defaultVerrazzanoName}
		defaultSecret := &v1.Secret{}
		if err := c.Get(context.TODO(), namespacedName, defaultSecret); err != nil {
			return err
		}
		if len(defaultSecret.Data[caCert]) < 1 {
			return errors.New(fmt.Sprintf("Failed, secret %s/%s does not have a value for %s", defaultSecretNamespace, defaultVerrazzanoName, caCert))

		}
		rancherCaSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: common.CattleSystem,
				Name:      rancherTLSSecretName,
			},
		}
		log.Printf("Copying default Verrazzano secret to Rancher namespace")
		if _, err := ctrl.CreateOrUpdate(context.TODO(), c, rancherCaSecret, func() error {
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
