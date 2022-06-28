// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package adapter

import (
	vzapi "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/authproxy"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/certmanager"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/externaldns"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/keycloak"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/oam"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/weblogic"
	"github.com/verrazzano/verrazzano/platform-operator/internal/config"
	"github.com/verrazzano/verrazzano/platform-operator/internal/vzconfig"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"os"
	"path"
	clipkg "sigs.k8s.io/controller-runtime/pkg/client"
)

const valuesYaml = "values.yaml"

func ApplyComponentAsModule(client clipkg.Client, vz *vzapi.Verrazzano, componentName string) error {
	adapter := componentAdapters[componentName]
	if adapter != nil {
		return adapter(vz).createOrUpdate(client)
	}
	return nil
}

func overridesData(fileName string) []byte {
	data, err := os.ReadFile(path.Join(config.GetHelmOverridesDir(), fileName))
	if err != nil {
		panic(err)
	}
	return data
}

// Adapter can be a part of the Component interface
var componentAdapters = map[string]func(*vzapi.Verrazzano) *componentAdapter{
	// oam adapter
	oam.ComponentName: func(vz *vzapi.Verrazzano) *componentAdapter {
		adapter := NewAdapter(vzconfig.IsOAMEnabled(vz))
		if adapter.IsEnabled {
			adapter.Name = oam.ComponentName
			adapter.Namespace = vz.Namespace
			adapter.ChartNamespace = oam.ComponentNamespace
			adapter.ChartPath = oam.ComponentName
			oam := vz.Spec.Components.OAM
			if oam != nil {
				adapter.InstallOverrides = oam.InstallOverrides
				override := vzapi.Overrides{
					Values: &apiextensionsv1.JSON{
						Raw: overridesData("oam-kubernetes-runtime-values.yaml"),
					},
				}
				adapter.InstallOverrides.ValueOverrides = append([]vzapi.Overrides{override}, oam.ValueOverrides...)
			}
		}
		return adapter
	},

	// external DNS adapter
	externaldns.ComponentName: func(vz *vzapi.Verrazzano) *componentAdapter {
		adapter := NewAdapter(vzconfig.IsExternalDNSEnabled(vz))
		if adapter.IsEnabled {
			adapter.Name = externaldns.ComponentName
			adapter.Namespace = vz.Namespace
			adapter.ChartNamespace = externaldns.ComponentNamespace
			adapter.ChartPath = externaldns.ComponentName
			dns := vz.Spec.Components.DNS
			if dns != nil {
				adapter.InstallOverrides = dns.InstallOverrides
				override := vzapi.Overrides{
					ConfigMapRef: &corev1.ConfigMapKeySelector{
						Key: valuesYaml,
						LocalObjectReference: corev1.LocalObjectReference{
							Name: externaldns.ConfigMapName,
						},
					},
				}
				adapter.InstallOverrides.ValueOverrides = append([]vzapi.Overrides{override}, dns.ValueOverrides...)
			}
		}
		return adapter
	},

	// cert manager adapter
	certmanager.ComponentName: func(vz *vzapi.Verrazzano) *componentAdapter {
		adapter := NewAdapter(vzconfig.IsCertManagerEnabled(vz))
		if adapter.IsEnabled {
			adapter.Name = certmanager.ComponentName
			adapter.Namespace = vz.Namespace
			adapter.ChartNamespace = certmanager.ComponentNamespace
			adapter.ChartPath = certmanager.ComponentName
			cm := vz.Spec.Components.CertManager
			if cm != nil {
				adapter.InstallOverrides = cm.InstallOverrides
				override := vzapi.Overrides{
					ConfigMapRef: &corev1.ConfigMapKeySelector{
						Key: valuesYaml,
						LocalObjectReference: corev1.LocalObjectReference{
							Name: certmanager.ConfigMapName,
						},
					},
				}
				adapter.InstallOverrides.ValueOverrides = append([]vzapi.Overrides{override}, cm.ValueOverrides...)
			}
		}
		return adapter
	},

	// Authproxy adapter
	authproxy.ComponentName: func(vz *vzapi.Verrazzano) *componentAdapter {
		adapter := NewAdapter(vzconfig.IsAuthProxyEnabled(vz))
		if adapter.IsEnabled {
			adapter.Name = authproxy.ComponentName
			adapter.Namespace = vz.Namespace
			adapter.ChartNamespace = authproxy.ComponentNamespace
			adapter.ChartPath = authproxy.ComponentName
		}
		return adapter
	},

	// Keycloak Adapter
	keycloak.ComponentName: func(vz *vzapi.Verrazzano) *componentAdapter {
		adapter := NewAdapter(vzconfig.IsKeycloakEnabled(vz))
		if adapter.IsEnabled {
			adapter.Name = keycloak.ComponentName
			adapter.Namespace = vz.Namespace
			adapter.ChartNamespace = keycloak.ComponentNamespace
			adapter.ChartPath = keycloak.ComponentName
			kc := vz.Spec.Components.Keycloak
			if kc != nil {
				adapter.InstallOverrides = kc.InstallOverrides
				override := vzapi.Overrides{
					ConfigMapRef: &corev1.ConfigMapKeySelector{
						Key: valuesYaml,
						LocalObjectReference: corev1.LocalObjectReference{
							Name: keycloak.ConfigMapName,
						},
					},
				}
				adapter.InstallOverrides.ValueOverrides = append([]vzapi.Overrides{override}, kc.ValueOverrides...)
			}
		}
		return adapter
	},

	// Weblogic Operator Adapter
	weblogic.ComponentName: func(vz *vzapi.Verrazzano) *componentAdapter {
		adapter := NewAdapter(vzconfig.IsWeblogicOperatorEnabled(vz))
		if adapter.IsEnabled {
			wko := vz.Spec.Components.WebLogicOperator
			adapter.Name = weblogic.ComponentName
			adapter.Namespace = vz.Namespace
			adapter.ChartNamespace = weblogic.ComponentNamespace
			adapter.ChartPath = weblogic.ComponentName
			if wko != nil {
				adapter.InstallOverrides = wko.InstallOverrides
				override := vzapi.Overrides{
					ConfigMapRef: &corev1.ConfigMapKeySelector{
						Key: valuesYaml,
						LocalObjectReference: corev1.LocalObjectReference{
							Name: weblogic.ConfigMapName,
						},
					},
				}
				adapter.InstallOverrides.ValueOverrides = append([]vzapi.Overrides{override}, wko.InstallOverrides.ValueOverrides...)
			}
		}
		return adapter
	},
}
