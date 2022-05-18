// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package metricstrait

import (
	"context"
	"fmt"

	promoperapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	vzapi "github.com/verrazzano/verrazzano/application-operator/apis/oam/v1alpha1"
	vzlog2 "github.com/verrazzano/verrazzano/pkg/log/vzlog"
	k8sapps "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const podMonitorSpec = `
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: ` + monitorNameHolder + `
  namespace: ` + namespaceHolder + `
  labels:
    name: ` + monitorNameHolder + `
spec:
  selector:
    matchLabels:
      app: "` + appNameHolder + `"
  podMetricsEndpoints:
  - port: http
    scheme: https
    tlsConfig:
      ca:
        secret:
          name: istio-certs
          key: root-cert.pem
      cert:
        secret:
          name: istio-certs
          key: cert-chain.pem
      keySecret:
        name: istio-certs
        key: key.pem
      insecureSkipVerify: true
    relabelings:
      - action: replace
        replacement: local
        sourceLabels: null
        targetLabel: ` + prometheusClusterNameLabel + `
      - action: keep
        regex: true;` + appNameHolder + `;` + compNameHolder + `
        sourceLabels:
        - __meta_kubernetes_pod_annotation_verrazzano_io_metricsEnabled` + portOrderHolder + `
        - __meta_kubernetes_pod_label_app_oam_dev_name
        - __meta_kubernetes_pod_label_app_oam_dev_component
      - action: replace
        regex: (.+)
        sourceLabels:
        - __meta_kubernetes_pod_annotation_verrazzano_io_metricsPath` + portOrderHolder + `
        targetLabel: __metrics_path__
      - action: replace
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
        sourceLabels:
        - __address__
        - __meta_kubernetes_pod_annotation_verrazzano_io_metricsPort` + portOrderHolder + `
        targetLabel: __address__
      - action: replace
        regex: (.*)
        replacement: $1
        sourceLabels:
        - __meta_kubernetes_namespace
        targetLabel: namespace
      - action: labelmap
        regex: __meta_kubernetes_pod_label_(.+)
      - action: replace
        sourceLabels:
        - __meta_kubernetes_pod_name
        targetLabel: pod_name
      - action: labeldrop
        regex: (controller_revision_hash)
      - action: replace
        regex: .*/(.*)$
        replacement: $1
        sourceLabels:
        - name
        targetLabel: webapp
`

// deleteOrUpdatePodMonitor reconciles a PodMonitor for the metrics trait.
// The PodMonitor for the trait is removed if needed.
func (r *Reconciler) deleteOrUpdatePodMonitor(ctx context.Context, trait *vzapi.MetricsTrait, rel vzapi.QualifiedResourceRelation, log vzlog2.VerrazzanoLogger) (vzapi.QualifiedResourceRelation, controllerutil.OperationResult, error) {
	deployment := &k8sapps.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: rel.Namespace, Name: rel.Name}, deployment)
	if err != nil {
		return rel, controllerutil.OperationResultNone, client.IgnoreNotFound(err)
	}
	return r.updatePodMonitor(ctx, trait, rel, log)
}

func (r *Reconciler) updatePodMonitor(ctx context.Context, trait *vzapi.MetricsTrait, rel vzapi.QualifiedResourceRelation, log vzlog2.VerrazzanoLogger) (vzapi.QualifiedResourceRelation, controllerutil.OperationResult, error) {
	podMonitor := promoperapi.PodMonitor{}
	name, err := createPodMonitorName(trait)
	if err != nil {
		return nil, controllerutil.OperationResultNone, err
	}
	namespace := getNamespaceFromObjectMetaOrDefault(trait.ObjectMeta)
	controllerutil.CreateOrUpdate(ctx, r.Client, &podMonitor, func() error {
		mutatePodMonitorFromTrait(&podMonitor, trait)
		return nil
	})
	err = r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &podMonitor)
	if err != nil {
	}
	// dummy return for now
	return vzapi.QualifiedResourceRelation{}, controllerutil.OperationResultNone, nil
}

func mutatePodMonitorFromTrait(podMonitor *promoperapi.PodMonitor, trait *vzapi.MetricsTrait, traitDefaults *vzapi.MetricsTraitSpec) {
	if podMonitor.ObjectMeta.Labels == nil {
		podMonitor.ObjectMeta.Labels = map[string]string{}
	}
	podMonitor.ObjectMeta.Labels["name"] = trait.Labels[appObjectMetaLabel]
	ports := getPortSpecs(trait, traitDefaults)
	for _, port := range ports {
		// TODO create a PodMetricsEndpoint for each port
		// TODO question: How do we currently get port specs from (e.g.) helidon metrics trait when
		//  the trait has no ports or path at all in its spec?
		//  getPortSpecs doesn't seem to do this. (ports = []vzapi.PortSpec{{Port: trait.Spec.Port, Path: trait.Spec.Path}})
		findPodMetricsEndpoint(port.Path)
		// podMonitor.Spec.PodMetricsEndpoints
	}
}

func createPodMonitorName(trait *vzapi.MetricsTrait) (string, error) {
	cluster := getClusterNameFromObjectMetaOrDefault(trait.ObjectMeta)
	app, found := trait.Labels[appObjectMetaLabel]
	if !found {
		return "", fmt.Errorf("metrics trait missing application name label")
	}
	comp, found := trait.Labels[compObjectMetaLabel]
	if !found {
		return "", fmt.Errorf("metrics trait missing component name label")
	}
	return fmt.Sprintf("%s-%s-%s", app, cluster, comp), nil
}
