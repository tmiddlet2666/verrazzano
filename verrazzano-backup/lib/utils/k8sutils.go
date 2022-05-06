// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	model "github.com/verrazzano/verrazzano/verrazzano-backup/lib/types"
	"go.uber.org/zap"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

type K8s interface {
	PopulateConnData(dclient dynamic.Interface, client client.Client, veleroNamespace, backupName, profile string, log *zap.SugaredLogger) (*model.ConnectionData, error)
	GetObjectStoreCreds(client client.Client, secretName, namespace, secretKey, profile string, log *zap.SugaredLogger) (*model.ObjectStoreSecret, error)
	GetBackup(client dynamic.Interface, veleroNamespace, backupName string, log *zap.SugaredLogger) (*model.VeleroBackup, error)
	GetBackupStorageLocation(client dynamic.Interface, veleroNamespace, bslName string, log *zap.SugaredLogger) (*model.VeleroBackupStorageLocation, error)
	ScaleDeployment(clientk client.Client, k8sclient *kubernetes.Clientset, labelSelector, namespace, deploymentName string, replicaCount int32, log *zap.SugaredLogger) error
}

type K8sImpl struct {
}

//PopulateConnData crestes the connection object thats used to communicate to object store
func (k *K8sImpl) PopulateConnData(dclient dynamic.Interface, client client.Client, veleroNamespace, backupName, profile string, log *zap.SugaredLogger) (*model.ConnectionData, error) {
	log.Infof("Populating connection data from backup '%v' in namespace '%s'", backupName, veleroNamespace)

	backup, err := k.GetBackup(dclient, veleroNamespace, backupName, log)
	if err != nil {
		return nil, err
	}

	if backup.Spec.StorageLocation == "default" {
		log.Infof("Default creds not supported. Custom credentaisl needs to be created before creating backup storage location")
		return nil, err
	}

	log.Infof("Detected velero backup storage location '%s' in namespace '%s' used by backup '%s'", backup.Spec.StorageLocation, veleroNamespace, backupName)
	bsl, err := k.GetBackupStorageLocation(dclient, veleroNamespace, backup.Spec.StorageLocation, log)
	if err != nil {
		return nil, err
	}

	secretData, err := k.GetObjectStoreCreds(client, bsl.Spec.Credential.Name, bsl.Metadata.Namespace, bsl.Spec.Credential.Key, profile, log)
	if err != nil {
		return nil, err
	}

	var conData model.ConnectionData
	conData.Secret = *secretData
	conData.RegionName = bsl.Spec.Config.Region
	conData.Endpoint = bsl.Spec.Config.S3URL
	conData.BucketName = bsl.Spec.ObjectStorage.Bucket
	conData.BackupName = backupName

	return &conData, nil

}

//GetObjectStoreCreds - Fetches credentials from Velero Backup object store location.
//This object will be pre-created before the execution of this hook
func (k *K8sImpl) GetObjectStoreCreds(client client.Client, secretName, namespace, secretKey, profile string, log *zap.SugaredLogger) (*model.ObjectStoreSecret, error) {
	secret := v1.Secret{}
	if err := client.Get(context.TODO(), crtclient.ObjectKey{Name: secretName, Namespace: namespace}, &secret); err != nil {
		log.Errorf("Failed to retrieve secret '%s' due to : %v", secretName, err)
		return nil, err
	}

	file, err := CreateTempFileWithData(secret.Data[secretKey])
	if err != nil {
		return nil, err
	}
	defer os.Remove(file)

	pathElements := strings.Split(file, "/")
	viper.SetConfigName(pathElements[len(pathElements)-1])
	viper.SetConfigType("ini")
	viper.AddConfigPath("/tmp/")
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var secretData model.ObjectStoreSecret
	secretData.SecretName = secretName
	secretData.SecretKey = secretKey

	accessKeyString := fmt.Sprintf("%s.aws_access_key_id", profile)
	secretData.ObjectAccessKey = fmt.Sprintf("%s", viper.Get(accessKeyString))
	secretAccessKeyString := fmt.Sprintf("%s.aws_secret_access_key", profile)
	secretData.ObjectSecretKey = fmt.Sprintf("%s", viper.Get(secretAccessKeyString))

	return &secretData, nil
}

//GetBackupStorageLocation - Retrieves the backup storage location from the backup storage location
func (k *K8sImpl) GetBackupStorageLocation(client dynamic.Interface, veleroNamespace, bslName string, log *zap.SugaredLogger) (*model.VeleroBackupStorageLocation, error) {
	log.Infof("Fetching velero backup storage location '%s' in namespace '%s'", bslName, veleroNamespace)
	gvr := schema.GroupVersionResource{
		Group:    "velero.io",
		Version:  "v1",
		Resource: "backupstoragelocations",
	}
	bslRecievd, err := client.Resource(gvr).Namespace(veleroNamespace).Get(context.Background(), bslName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if bslRecievd == nil {
		log.Infof("No velero backup storage location in namespace '%s' was detected", veleroNamespace)
		return nil, err
	}

	var bsl model.VeleroBackupStorageLocation
	bdata, err := json.Marshal(bslRecievd)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bdata, &bsl)
	if err != nil {
		return nil, err
	}
	return &bsl, nil
}

//GetBackup - Retrives velero backups in the cluster
func (k *K8sImpl) GetBackup(client dynamic.Interface, veleroNamespace, backupName string, log *zap.SugaredLogger) (*model.VeleroBackup, error) {
	log.Infof("Fetching velero backup '%s' in namespace '%s'", backupName, veleroNamespace)
	gvr := schema.GroupVersionResource{
		Group:    "velero.io",
		Version:  "v1",
		Resource: "backups",
	}
	backupFetched, err := client.Resource(gvr).Namespace(veleroNamespace).Get(context.Background(), backupName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if backupFetched == nil {
		log.Infof("No velero backup in namespace '%s' was detected", veleroNamespace)
		return nil, err
	}

	var backup model.VeleroBackup
	bdata, err := json.Marshal(backupFetched)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bdata, &backup)
	if err != nil {
		return nil, err
	}
	return &backup, nil
}

//ScaleDeployment is used to scale deployment to specific replica count
// labelselectors,namespace, deploymentName are used to identify deployments and specific pods associated with them
func (k *K8sImpl) ScaleDeployment(clientk client.Client, k8sclient *kubernetes.Clientset, labelSelector, namespace, deploymentName string, replicaCount int32, log *zap.SugaredLogger) error {
	log.Infof("Scale deployment '%s' in namespace '%s", deploymentName, namespace)
	depPatch := apps.Deployment{}
	if err := clientk.Get(context.TODO(), types.NamespacedName{Name: deploymentName, Namespace: namespace}, &depPatch); err != nil {
		return err
	}
	currentValue := *depPatch.Spec.Replicas
	desiredValue := replicaCount

	if desiredValue == currentValue {
		log.Infof("Deployment scaling skipped as desired replicas is same as current replicas")
		return nil
	}

	mergeFromDep := client.MergeFrom(depPatch.DeepCopy())
	depPatch.Spec.Replicas = &replicaCount
	if err := clientk.Patch(context.TODO(), &depPatch, mergeFromDep); err != nil {
		log.Error("Unable to patch !!")
		return err
	}

	done := false
	listOptions := metav1.ListOptions{LabelSelector: labelSelector}
	var podStateCondition []bool

	for !done {
		pods, err := k8sclient.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
		if err != nil {
			return err
		}

		//Scale up
		if desiredValue > currentValue {
			log.Info("Scaling up ...")
			// There could be multiple pods in a deployment
			for _, item := range pods.Items {
				if item.Status.Phase == "Running" {
					podStateCondition = append(podStateCondition, true)
				}
			}

			if int32(len(pods.Items)) == desiredValue && int32(len(podStateCondition)) == desiredValue {
				// when all running pods is equal to desired input
				// exit the check loop
				done = true
			} else {
				// otherwise retry and keep monitoring the pod status
				log.Info("Waiting for 30 seconds for all pods to come up.")
				time.Sleep(time.Second * 30)
			}

		}

		// scale down
		if desiredValue < currentValue {
			log.Info("Scaling down ..")
			if int32(len(pods.Items)) != desiredValue {
				log.Info("Waiting for 30 seconds for all  pods to go down.")
				time.Sleep(time.Second * 30)
			} else {
				done = true
			}
		}

	}

	log.Infof("Successfully scaled deployment '%s' in namespace '%s' from '%v' to '%v' replicas ", deploymentName, namespace, currentValue, replicaCount)
	return nil

}
