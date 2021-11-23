// Copyright (c) 2021, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package rancher

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestCreateAdminSecretIfNotExists(t *testing.T) {
	log := getTestLogger(t)
	ff := func(a ...string) (string, string, error) {
		return "", "", nil
	}

	podList := createRancherPodList()
	adminSecret := createAdminSecret()

	var tests = []struct {
		testName string
		c        client.Client
		f        func(a ...string) (string, string, error)
		isErr    bool
	}{
		{
			"should skip secret creation when secret is present",
			fake.NewFakeClientWithScheme(getScheme(), &adminSecret),
			ff,
			false,
		},
		{
			"should be able to reset the admin password",
			fake.NewFakeClientWithScheme(getScheme(), &podList),
			func(a ...string) (string, string, error) {
				return "password", "", nil
			},
			false,
		},
		{
			"should fail when resetting admin password fails",
			fake.NewFakeClientWithScheme(getScheme(), &podList),
			func(a ...string) (string, string, error) {
				return "", "", errors.New("something bad happened")
			},
			true,
		},
		{
			"should fail when no rancher pod is available",
			fake.NewFakeClientWithScheme(getScheme()),
			ff,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			setBashFunc(tt.f)
			err := createAdminSecretIfNotExists(log, tt.c)
			if tt.isErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
