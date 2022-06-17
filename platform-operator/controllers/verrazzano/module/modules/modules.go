package modules

import (
	"context"
	modulesv1alpha1 "github.com/verrazzano/verrazzano/platform-operator/apis/modules/v1alpha1"
	clipkg "sigs.k8s.io/controller-runtime/pkg/client"
)

func UpdateModule(client clipkg.Client, module *modulesv1alpha1.Module) error {
	return client.Update(context.TODO(), module)
}
