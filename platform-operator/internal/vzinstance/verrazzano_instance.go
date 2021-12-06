// Copyright (C) 2021, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package vzinstance

import (
	"context"
	"fmt"
	"github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/constants"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"

)

const systemNamespace = "verrazzano-system"

// GetInstanceInfo returns the instance info for the local install.
func GetInstanceInfo(client client.Client, cr *v1alpha1.Verrazzano) *v1alpha1.InstanceInfo {

	ingressList := &networkingv1.IngressList{}
	err := client.List(context.TODO(), ingressList)
	if err != nil {
		zap.S().Errorf("Error listing ingresses: %v", err)
		return nil
	}
	if len(ingressList.Items) == 0 {
		zap.S().Debugf("No ingresses found, unable to build instance info")
		return nil
	}
	svcList := &corev1.ServiceList{}
	err = client.List(context.TODO(), svcList)
	if err != nil {
		zap.S().Errorf("Error listing ingresses: %v", err)
		return nil
	}
	if len(svcList.Items) == 0 {
		zap.S().Debugf("No services found, unable to build instance info")
		return nil
	}

	svcNg:= findService(svcList.Items,"ingress-nginx","ingress-controller-ingress-nginx-controller")
	fmt.Printf("+++ ABMITRA svcType = %v ++++\n",svcNg.Spec.Type)


	var nodePort int32
	if svcNg.Spec.Type == "NodePort" {
		svc := findService(svcList.Items,"ingress-nginx","ingress-controller-ingress-nginx-controller")
		for _,port := range svc.Spec.Ports {
			if port.Port == 443 {
				nodePort = port.NodePort
			}
		}
	}

	// Console ingress always exist. Only show console URL if the console was enabled during install.
	var consoleURL *string
	if cr.Spec.Components.Console == nil || *cr.Spec.Components.Console.Enabled {
		consoleURL = getSystemIngressURL(ingressList.Items, systemNamespace, constants.VzConsoleIngress,nodePort)
	} else {
		consoleURL = nil
	}
	instanceInfo := &v1alpha1.InstanceInfo{
		ConsoleURL:    consoleURL,
		RancherURL:    getSystemIngressURL(ingressList.Items, "cattle-system", "rancher",nodePort),
		KeyCloakURL:   getSystemIngressURL(ingressList.Items, "keycloak", "keycloak",nodePort),
		ElasticURL:    getSystemIngressURL(ingressList.Items, systemNamespace, "vmi-system-es-ingest",nodePort),
		KibanaURL:     getSystemIngressURL(ingressList.Items, systemNamespace, "vmi-system-kibana",nodePort),
		GrafanaURL:    getSystemIngressURL(ingressList.Items, systemNamespace, "vmi-system-grafana",nodePort),
		PrometheusURL: getSystemIngressURL(ingressList.Items, systemNamespace, "vmi-system-prometheus",nodePort),
		KialiURL:      getSystemIngressURL(ingressList.Items, systemNamespace, "vmi-system-kiali",nodePort),
	}
	return instanceInfo
}

func getSystemIngressURL(ingresses []networkingv1.Ingress, namespace string, name string,nodePort int32) *string {
	var ingress = findIngress(ingresses, namespace, name)
	if ingress == nil {
		zap.S().Debugf("No ingress found for %s/%s", namespace, name)
		return nil
	}
	var url string
	if nodePort > 0  {
		nodePortString := strconv.Itoa(int(nodePort))
		url = fmt.Sprintf("https://%s:%s", ingress.Spec.Rules[0].Host,nodePortString)
	} else {
		url = fmt.Sprintf("https://%s", ingress.Spec.Rules[0].Host)
	}
	return &url
}

func findIngress(ingresses []networkingv1.Ingress, namespace string, name string) *networkingv1.Ingress {
	for _, ingress := range ingresses {
		if ingress.Name == name && ingress.Namespace == namespace {
			return &ingress
		}
	}
	return nil
}

func findService(services []corev1.Service, namespace string, name string) *corev1.Service{
	for _, service := range services {
		if service.Name == name && service.Namespace == namespace {
			return &service
		}
	}
	return nil
}