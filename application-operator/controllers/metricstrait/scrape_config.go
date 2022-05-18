// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package metricstrait

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Jeffail/gabs/v2"
	vzapi "github.com/verrazzano/verrazzano/application-operator/apis/oam/v1alpha1"
	"github.com/verrazzano/verrazzano/application-operator/controllers/clusters"
	vznav "github.com/verrazzano/verrazzano/application-operator/controllers/navigation"
	vzlog2 "github.com/verrazzano/verrazzano/pkg/log/vzlog"
	k8sapps "k8s.io/api/apps/v1"
	k8score "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// prometheusScrapeConfigTemplate configuration for general Prometheus scrape target template
// Used to add new scrape config to a Prometheus configmap
const prometheusScrapeConfigTemplate = `job_name: ##JOB_NAME##
##SSL_PROTOCOL##
kubernetes_sd_configs:
- role: pod
  namespaces:
    names:
    - ##NAMESPACE##
relabel_configs:
- action: replace
  source_labels: null
  target_label: ` + prometheusClusterNameLabel + `
  replacement: ##VERRAZZANO_CLUSTER_NAME##
- action: keep
  source_labels: [__meta_kubernetes_pod_annotation_verrazzano_io_metricsEnabled##PORT_ORDER##,__meta_kubernetes_pod_label_app_oam_dev_name,__meta_kubernetes_pod_label_app_oam_dev_component]
  regex: true;##APP_NAME##;##COMP_NAME##
- action: replace
  source_labels: [__meta_kubernetes_pod_annotation_verrazzano_io_metricsPath##PORT_ORDER##]
  target_label: __metrics_path__
  regex: (.+)
- action: replace
  source_labels: [__address__, __meta_kubernetes_pod_annotation_verrazzano_io_metricsPort##PORT_ORDER##]
  target_label: __address__
  regex: ([^:]+)(?::\d+)?;(\d+)
  replacement: $1:$2
- action: replace
  source_labels: [__meta_kubernetes_namespace]
  target_label: namespace
  regex: (.*)
  replacement: $1
- action: labelmap
  regex: __meta_kubernetes_pod_label_(.+)
- action: replace
  source_labels: [__meta_kubernetes_pod_name]
  target_label: pod_name
- action: labeldrop
  regex: '(controller_revision_hash)'
- action: replace
  source_labels: [name]
  target_label: webapp
  regex: '.*/(.*)$'
  replacement: $1
`

// prometheusWLSScrapeConfigTemplate configuration for WebLogic Prometheus scrape target template
// Used to add new WebLogic scrape config to a Prometheus configmap
const prometheusWLSScrapeConfigTemplate = `job_name: ##JOB_NAME##
##SSL_PROTOCOL##
kubernetes_sd_configs:
- role: pod
  namespaces:
    names:
    - ##NAMESPACE##
relabel_configs:
- action: replace
  source_labels: null
  target_label: ` + prometheusClusterNameLabel + `
  replacement: ##VERRAZZANO_CLUSTER_NAME##
- action: keep
  source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape,__meta_kubernetes_pod_label_app_oam_dev_name,__meta_kubernetes_pod_label_app_oam_dev_component]
  regex: true;##APP_NAME##;##COMP_NAME##
- action: replace
  source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
  target_label: __metrics_path__
  regex: (.+)
- action: replace
  source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
  target_label: __address__
  regex: ([^:]+)(?::\d+)?;(\d+)
  replacement: $1:$2
- action: replace
  source_labels: [__meta_kubernetes_namespace]
  target_label: namespace
  regex: (.*)
  replacement: $1
- action: labelmap
  regex: __meta_kubernetes_pod_label_(.+)
- action: replace
  source_labels: [__meta_kubernetes_pod_name]
  target_label: pod_name
- action: labeldrop
  regex: '(controller_revision_hash)'
- action: replace
  source_labels: [name]
  target_label: webapp
  regex: '.*/(.*)$'
  replacement: $1
`

// deleteOrUpdateScraperConfigMap cleans up a scraper (i.e. Prometheus) configmap.
// The scraper config for the trait is removed if present.
func (r *Reconciler) deleteOrUpdateScraperConfigMap(ctx context.Context, trait *vzapi.MetricsTrait, rel vzapi.QualifiedResourceRelation, log vzlog2.VerrazzanoLogger) (vzapi.QualifiedResourceRelation, controllerutil.OperationResult, error) {
	deployment := &k8sapps.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: rel.Namespace, Name: rel.Name}, deployment)
	if err != nil {
		return rel, controllerutil.OperationResultNone, client.IgnoreNotFound(err)
	}
	return r.updatePrometheusScraperConfigMap(ctx, trait, nil, nil, deployment, log)
}

// updatePrometheusScraperConfigMap updates the Prometheus scraper configmap.
// This updates only the scrape_configs section of the Prometheus configmap.
// Only the rules for the provided trait will be affected.
// trait - The trait to update scrape_config rules for.
// traitDefaults - Default to use for values not provided in the trait.
// deployment - The Prometheus deployment.
func (r *Reconciler) updatePrometheusScraperConfigMap(ctx context.Context, trait *vzapi.MetricsTrait, workload *unstructured.Unstructured, traitDefaults *vzapi.MetricsTraitSpec, deployment *k8sapps.Deployment, log vzlog2.VerrazzanoLogger) (vzapi.QualifiedResourceRelation, controllerutil.OperationResult, error) {
	rel := vzapi.QualifiedResourceRelation{APIVersion: deployment.APIVersion, Kind: deployment.Kind, Name: deployment.Name, Namespace: deployment.Namespace, Role: scraperRole}

	// Fetch the secret by name if it is provided in either the trait or the trait defaults.
	secret, err := r.fetchSourceCredentialsSecretIfRequired(ctx, trait, traitDefaults, workload)
	if err != nil {
		return rel, controllerutil.OperationResultNone, err
	}

	configmapName, err := r.findPrometheusScrapeConfigMapNameFromDeployment(deployment, log)
	if err != nil {
		return rel, controllerutil.OperationResultNone, err
	}

	configmap := &k8score.ConfigMap{}
	err = r.Get(ctx, client.ObjectKey{Namespace: deployment.Namespace, Name: configmapName}, configmap)
	if err != nil {
		// Don't create the config map if it doesn't already exist - that is the sole responsibility of
		// the Verrazzano Monitoring Operator
		return rel, controllerutil.OperationResultNone, client.IgnoreNotFound(err)
	}

	existingConfigmap := configmap.DeepCopyObject()

	if configmap.CreationTimestamp.IsZero() {
		log.Debugf("Create Prometheus configmap %s", vznav.GetNamespacedNameFromObjectMeta(configmap.ObjectMeta))
	} else {
		log.Debugf("Update Prometheus configmap %s", vznav.GetNamespacedNameFromObjectMeta(configmap.ObjectMeta))
	}
	yamlStr, exists := configmap.Data[prometheusConfigKey]
	if !exists {
		yamlStr = ""
	}
	prometheusConf, err := parseYAMLString(yamlStr)
	if err != nil {
		return rel, controllerutil.OperationResultNone, err
	}
	prometheusConf, err = mutatePrometheusScrapeConfig(ctx, trait, traitDefaults, prometheusConf, secret, workload, r.Client)
	if err != nil {
		return rel, controllerutil.OperationResultNone, err
	}
	yamlStr, err = writeYAMLString(prometheusConf)
	if err != nil {
		return rel, controllerutil.OperationResultNone, err
	}
	if configmap.Data == nil {
		configmap.Data = map[string]string{}
	}
	configmap.Data[prometheusConfigKey] = yamlStr

	// compare and don't update if unchanged
	if equality.Semantic.DeepEqual(existingConfigmap, configmap) {
		return rel, controllerutil.OperationResultNone, nil
	}

	err = r.Update(ctx, configmap)
	// If the Prometheus configmap was updated, the VMI Prometheus has ConfigReloader sidecar to signal Prometheus to reload config
	if err != nil {
		return rel, controllerutil.OperationResultNone, err
	}
	return rel, controllerutil.OperationResultUpdated, nil
}

// findPrometheusScrapeConfigMapNameFromDeployment finds the Prometheus configmap name from the Prometheus deployment.
func (r *Reconciler) findPrometheusScrapeConfigMapNameFromDeployment(deployment *k8sapps.Deployment, log vzlog2.VerrazzanoLogger) (string, error) {
	volumes := deployment.Spec.Template.Spec.Volumes
	for _, volume := range volumes {
		if volume.Name == "config-volume" && volume.ConfigMap != nil && len(volume.ConfigMap.Name) > 0 {
			name := volume.ConfigMap.Name
			log.Debugf("Found Prometheus configmap name %s", name)
			return name, nil
		}
	}
	return "", fmt.Errorf("failed to find Prometheus configmap name from deployment %s", vznav.GetNamespacedNameFromObjectMeta(deployment.ObjectMeta))
}

// mutatePrometheusScrapeConfig mutates the Prometheus scrape configuration.
// Scrap configuration rules will be added, updated, deleted depending on the state of the trait.
func mutatePrometheusScrapeConfig(ctx context.Context, trait *vzapi.MetricsTrait, traitDefaults *vzapi.MetricsTraitSpec, prometheusScrapeConfig *gabs.Container, secret *k8score.Secret, workload *unstructured.Unstructured, c client.Client) (*gabs.Container, error) {
	ports := getPortSpecs(trait, traitDefaults)

	for i := range ports {
		oldScrapeConfigs := prometheusScrapeConfig.Search(prometheusScrapeConfigsLabel).Children()
		prometheusScrapeConfig.Array(prometheusScrapeConfigsLabel) // zero out the array of scrape configs
		newScrapeJob, newScrapeConfig, err := createScrapeConfigFromTrait(ctx, trait, i, secret, workload, c)
		if err != nil {
			return prometheusScrapeConfig, err
		}
		existingReplaced := false
		for _, oldScrapeConfig := range oldScrapeConfigs {
			oldScrapeJob := oldScrapeConfig.Search(prometheusJobNameLabel).Data()
			if newScrapeJob == oldScrapeJob {
				// If the scrape config should be removed then skip adding it to the result slice.
				// This will occur in three situations.
				// 1. The trait is being deleted.
				// 2. The trait scraper has been changed and the old scrape config is being updated.
				//    In this case the traitDefaults and newScrapeConfig will be nil.
				// 3. The trait is being disabled.
				if trait.DeletionTimestamp.IsZero() && traitDefaults != nil && newScrapeConfig != nil && isEnabled(trait) {
					prometheusScrapeConfig.ArrayAppendP(newScrapeConfig.Data(), prometheusScrapeConfigsLabel)
				}
				existingReplaced = true
			} else {
				prometheusScrapeConfig.ArrayAppendP(oldScrapeConfig.Data(), prometheusScrapeConfigsLabel)
			}
		}
		// If an existing config was not replaced and there is new config (i.e. newScrapeConfig != nil) then add the new config.
		if !existingReplaced && newScrapeConfig != nil {
			prometheusScrapeConfig.ArrayAppendP(newScrapeConfig.Data(), prometheusScrapeConfigsLabel)
		}
	}
	return prometheusScrapeConfig, nil
}

// createScrapeConfigFromTrait creates Prometheus scrape config for a trait.
// This populates the Prometheus scrape config template.
// The job name is returned.
// The YAML container populated from the Prometheus scrape config template is returned.
func createScrapeConfigFromTrait(ctx context.Context, trait *vzapi.MetricsTrait, portIncrement int, secret *k8score.Secret, workload *unstructured.Unstructured, c client.Client) (string, *gabs.Container, error) {

	// TODO: see if we can create a scrape job per port within this method. change name to createScrapeConfigsFromTrait
	job, err := createPrometheusScrapeConfigMapJobName(trait, portIncrement)
	if err != nil {
		return "", nil, err
	}

	// If the metricsTrait is being disabled then return nil for the config
	if !isEnabled(trait) {
		return job, nil, nil
	}

	// If workload is nil then the trait is being deleted so no config is required
	if workload != nil {
		// Populate the Prometheus scrape config template
		portOrderStr := ""
		if portIncrement > 0 {
			portOrderStr = strconv.Itoa(portIncrement)
		}
		context := map[string]string{
			appNameHolder:       trait.Labels[appObjectMetaLabel],
			compNameHolder:      trait.Labels[compObjectMetaLabel],
			jobNameHolder:       job,
			portOrderHolder:     portOrderStr,
			namespaceHolder:     trait.Namespace,
			sslProtocolHolder:   httpProtocol,
			vzClusterNameHolder: clusters.GetClusterName(ctx, c)}

		var configTemplate string
		https, err := useHTTPSForScrapeTarget(ctx, c, trait)
		if err != nil {
			return "", nil, err
		}

		if https {
			context[sslProtocolHolder] = httpsProtocol
		}
		configTemplate = prometheusScrapeConfigTemplate
		apiVerKind, err := vznav.GetAPIVersionKindOfUnstructured(workload)
		if err != nil {
			return "", nil, err
		}
		// Match any version of APIVersion=weblogic.oracle and Kind=Domain
		if matched, _ := regexp.MatchString("^weblogic.oracle/.*\\.Domain$", apiVerKind); matched {
			configTemplate = prometheusWLSScrapeConfigTemplate
		}

		// Populate the Prometheus scrape config template
		template := mergeTemplateWithContext(configTemplate, context)

		// Parse the populate the Prometheus scrape config template.
		config, err := parseYAMLString(template)
		if err != nil {
			return job, nil, fmt.Errorf("failed to parse built-in Prometheus scrape config template: %w", err)
		}
		// Add basic auth credentials if provided
		if secret != nil {
			username, secretFound := secret.Data["username"]
			if secretFound {
				config.Set(string(username), basicAuthLabel, basicAuthUsernameLabel)
			}
			password, passwordFound := secret.Data["password"]
			if passwordFound {
				config.Set(string(password), basicAuthLabel, basicPathPasswordLabel)
			}
		}
		return job, config, nil
	}

	// If the trait is being deleted (i.e. workload==nil) then no config is required.
	return job, nil, nil
}

// createPrometheusScrapeConfigMapJobName creates a Prometheus scrape configmap job name from a trait.
// Format is {oam_app}_{cluster}_{namespace}_{oam_comp}
func createPrometheusScrapeConfigMapJobName(trait *vzapi.MetricsTrait, portNum int) (string, error) {
	cluster := getClusterNameFromObjectMetaOrDefault(trait.ObjectMeta)
	namespace := getNamespaceFromObjectMetaOrDefault(trait.ObjectMeta)
	app, found := trait.Labels[appObjectMetaLabel]
	if !found {
		return "", fmt.Errorf("metrics trait missing application name label")
	}
	comp, found := trait.Labels[compObjectMetaLabel]
	if !found {
		return "", fmt.Errorf("metrics trait missing component name label")
	}
	portStr := ""
	if portNum > 0 {
		portStr = fmt.Sprintf("_%d", portNum)
	}
	return fmt.Sprintf("%s_%s_%s_%s%s", app, cluster, namespace, comp, portStr), nil
}
