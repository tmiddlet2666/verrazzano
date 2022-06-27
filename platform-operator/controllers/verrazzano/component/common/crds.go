// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package common

import (
	"github.com/verrazzano/verrazzano/platform-operator/internal/config"
	"path/filepath"

	"github.com/verrazzano/verrazzano/pkg/k8sutil"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
)

func ApplyCRDYaml(ctx spi.ComponentContext, helmChartsDir string) error {
	path := filepath.Join(helmChartsDir, "/crds")
	yamlApplier := k8sutil.NewYAMLApplier(ctx.Client(), "")
	ctx.Log().Oncef("Applying yaml for crds in %s", path)
	return yamlApplier.ApplyD(path)
}

func ApplyOverride(ctx spi.ComponentContext, overrideFile string) error {
	yamlApplier := k8sutil.NewYAMLApplier(ctx.Client(), ctx.EffectiveCR().Namespace)
	path := filepath.Join(config.GetHelmOverridesDir(), overrideFile)
	ctx.Log().Oncef("Applying override objects in %s", path)
	return yamlApplier.ApplyF(path)
}
