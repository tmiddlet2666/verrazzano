// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package authproxy

import (
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/common"
	"io/fs"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verrazzano/verrazzano/pkg/bom"
	globalconst "github.com/verrazzano/verrazzano/pkg/constants"
	vzos "github.com/verrazzano/verrazzano/pkg/os"
	vzapi "github.com/verrazzano/verrazzano/platform-operator/apis/verrazzano/v1alpha1"
	vpoconst "github.com/verrazzano/verrazzano/platform-operator/constants"
	"github.com/verrazzano/verrazzano/platform-operator/controllers/verrazzano/component/spi"
	"github.com/verrazzano/verrazzano/platform-operator/internal/config"
	istioclinet "istio.io/client-go/pkg/apis/networking/v1alpha3"
	istioclisec "istio.io/client-go/pkg/apis/security/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/yaml"
)

var testScheme = runtime.NewScheme()

const (
	profileDir      = "../../../../manifests/profiles"
	testBomFilePath = "../../testdata/test_bom.json"
)

func init() {
	_ = clientgoscheme.AddToScheme(testScheme)
	_ = vzapi.AddToScheme(testScheme)
	_ = istioclinet.AddToScheme(testScheme)
	_ = istioclisec.AddToScheme(testScheme)
}

// TestIsAuthProxyReady tests the isAuthProxyReady call
// GIVEN a AuthProxy component
//  WHEN I call isAuthProxyReady when all requirements are met
//  THEN true or false is returned
func TestIsAuthProxyReady(t *testing.T) {
	tests := []struct {
		name       string
		client     client.Client
		expectTrue bool
	}{
		{
			name: "Test IsReady when AuthProxy is successfully deployed",
			client: fake.NewClientBuilder().WithScheme(testScheme).WithObjects(
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ComponentNamespace,
						Name:      ComponentName,
						Labels:    map[string]string{"app": ComponentName},
					},
					Status: appsv1.DeploymentStatus{
						AvailableReplicas: 1,
						Replicas:          1,
						UpdatedReplicas:   1,
					},
				}).Build(),
			expectTrue: true,
		},
		{
			name: "Test IsReady when AuthProxy deployment is not ready",
			client: fake.NewClientBuilder().WithScheme(testScheme).WithObjects(
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ComponentNamespace,
						Name:      ComponentName,
					},
					Status: appsv1.DeploymentStatus{
						AvailableReplicas: 1,
						Replicas:          1,
						UpdatedReplicas:   0,
					},
				}).Build(),
			expectTrue: false,
		},
		{
			name:       "Test IsReady when AuthProxy deployment does not exist",
			client:     fake.NewClientBuilder().WithScheme(testScheme).Build(),
			expectTrue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := spi.NewFakeContext(tt.client, &vzapi.Verrazzano{}, false)
			if tt.expectTrue {
				assert.True(t, isAuthProxyReady(ctx))
			} else {
				assert.False(t, isAuthProxyReady(ctx))
			}
		})
	}
}

// TestAppendOverrides tests the AppendOverrides function
// GIVEN a call to AppendOverrides
//  WHEN I call with a ComponentContext with different profiles and overrides
//  THEN the correct overrides file is generated
//
// For each test case a Verrazzano custom resource with different overrides
// is passed into AppendOverrides.  A overrides file is generated by AppendOverrides.
// The test compares the generated and expected overrides files.
func TestAppendOverrides(t *testing.T) {
	config.SetDefaultBomFilePath(testBomFilePath)
	defer func() {
		config.SetDefaultBomFilePath("")
	}()
	tests := []struct {
		name         string
		description  string
		expectedYAML string
		actualCR     string
		numKeyValues int
		expectedErr  error
	}{
		{
			name:         "DefaultConfig",
			description:  "Test default configuration of AuthProxy with no overrides",
			expectedYAML: "testdata/noOverrideValues.yaml",
			actualCR:     "testdata/noOverrideVz.yaml",
			numKeyValues: 1,
			expectedErr:  nil,
		},
		{
			name:         "OverrideReplicas",
			description:  "Test override of replica count",
			expectedYAML: "testdata/replicasOverrideValues.yaml",
			actualCR:     "testdata/replicasOverrideVz.yaml",
			numKeyValues: 1,
			expectedErr:  nil,
		},
		{
			name:         "OverrideAffinity",
			description:  "Test override of affinity configuration for AuthProxy",
			expectedYAML: "testdata/affinityOverrideValues.yaml",
			actualCR:     "testdata/affinityOverrideVz.yaml",
			numKeyValues: 1,
			expectedErr:  nil,
		},
		{
			name:         "OverrideDNSWildcardDomain",
			description:  "Test overriding DNS wildcard domain",
			expectedYAML: "testdata/dnsWildcardDomainOverrideValues.yaml",
			actualCR:     "testdata/dnsWildcardDomainOverrideVz.yaml",
			numKeyValues: 1,
			expectedErr:  nil,
		},
		{
			name:         "DisableAuthProxy",
			description:  "Test overriding AuthProxy to be disabled",
			expectedYAML: "testdata/enabledOverrideValues.yaml",
			actualCR:     "testdata/enabledOverrideVz.yaml",
			numKeyValues: 1,
			expectedErr:  nil,
		},
	}
	defer resetWriteFileFunc()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			asserts := assert.New(t)
			t.Log(test.description)

			// Read the Verrazzano CR into a struct
			testCR := vzapi.Verrazzano{}
			yamlFile, err := ioutil.ReadFile(test.actualCR)
			asserts.NoError(err)
			err = yaml.Unmarshal(yamlFile, &testCR)
			asserts.NoError(err)

			fakeClient := createFakeClientWithIngress()
			fakeContext := spi.NewFakeContext(fakeClient, &testCR, false, profileDir)

			writeFileFunc = func(filename string, data []byte, perm fs.FileMode) error {
				if test.expectedErr != nil {
					return test.expectedErr
				}
				if err := ioutil.WriteFile(filename, data, perm); err != nil {
					asserts.Failf("Failure writing file %s: %s", filename, err)
					return err
				}
				asserts.FileExists(filename)

				// Unmarshal the actual generated helm values from code under test
				actualValues := authProxyValues{}
				err := yaml.Unmarshal(data, &actualValues)
				asserts.NoError(err)

				// read in the expected results' data from a file and unmarshal it into a values object
				expectedData, err := ioutil.ReadFile(test.expectedYAML)
				asserts.NoError(err, "Error reading expected values yaml file %s", test.expectedYAML)
				expectedValues := authProxyValues{}
				err = yaml.Unmarshal(expectedData, &expectedValues)
				asserts.NoError(err)

				// Compare the actual and expected values objects
				asserts.Equal(expectedValues, actualValues)
				return nil
			}

			var kvs []bom.KeyValue
			kvs, err = AppendOverrides(fakeContext, "", "", "", kvs)
			if test.expectedErr != nil {
				asserts.Error(err)
				asserts.Equal([]bom.KeyValue{}, kvs)
				return
			}
			asserts.NoError(err)
			asserts.Equal(test.numKeyValues, len(kvs))

			// Check Temp file
			asserts.True(kvs[0].IsFile, "Expected generated AuthProxy overrides first in list of helm args")
			tempFilePath := kvs[0].Value
			_, err = os.Stat(tempFilePath)
			asserts.NoError(err, "Unexpected error checking for temp file %s: %s", tempFilePath, err)
			cleanTempFiles(fakeContext)
		})
	}
	// Verify temp files are deleted
	files, err := ioutil.ReadDir(os.TempDir())
	assert.NoError(t, err, "Error reading temp dir to verify file cleanup")
	for _, file := range files {
		assert.False(t,
			strings.HasPrefix(file.Name(), tmpFilePrefix) && strings.HasSuffix(file.Name(), ".yaml"),
			"Found unexpected temp file remaining: %s", file.Name())
	}

}

// TestRemoveResourcePolicyAnnotation tests the removeResourcePolicyAnnotation function
// GIVEN a call to removeResourcePolicyAnnotation
//  WHEN I call with a object that is annotated with the resource policy annotation
//  THEN the annotation is removed
func TestRemoveResourcePolicyAnnotation(t *testing.T) {
	namespacedName := types.NamespacedName{
		Name:      ComponentName,
		Namespace: ComponentNamespace,
	}
	obj := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ComponentName,
			Namespace: ComponentNamespace,
			Annotations: map[string]string{"meta.helm.sh/release-name": ComponentName, "meta.helm.sh/release-namespace": ComponentNamespace,
				"helm.sh/resource-policy": "keep"},
		},
	}

	c := fake.NewClientBuilder().WithScheme(testScheme).WithObjects(obj).Build()
	res, err := common.RemoveResourcePolicyAnnotation(c, obj, namespacedName)
	assert.NoError(t, err)
	assert.Equal(t, ComponentName, res.GetAnnotations()["meta.helm.sh/release-name"])
	assert.Equal(t, globalconst.VerrazzanoSystemNamespace, res.GetAnnotations()["meta.helm.sh/release-namespace"])
	_, ok := res.GetAnnotations()["helm.sh/resource-policy"]
	assert.False(t, ok)
}

func createFakeClientWithIngress() client.Client {
	fakeClient := fake.NewClientBuilder().WithScheme(testScheme).WithObjects(
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{Name: vpoconst.NGINXControllerServiceName, Namespace: globalconst.IngressNamespace},
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeLoadBalancer,
			},
			Status: corev1.ServiceStatus{
				LoadBalancer: corev1.LoadBalancerStatus{
					Ingress: []corev1.LoadBalancerIngress{
						{IP: "11.22.33.44"},
					},
				},
			},
		},
	).Build()
	return fakeClient
}

//cleanTempFiles - Clean up the override temp files in the temp dir
func cleanTempFiles(ctx spi.ComponentContext) {
	if err := vzos.RemoveTempFiles(ctx.Log().GetZapLogger(), tmpFileCleanPattern); err != nil {
		ctx.Log().Errorf("Failed deleting temp files: %v", err)
	}
}
