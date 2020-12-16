package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/util/homedir"
)

// TestGetKubeconfigEnvVar tests getKubeconfig when
// KUBECONFIG environment variable is set
func TestGetKubeconfigEnvVar(t *testing.T) {
	// dont overwrite the kubeconfig path if it exists
	var tmp string
	if os.Getenv("KUBECONFIG") != "" {
		os.Unsetenv("KUBECONFIG")
		// re-write upon exit
		defer os.Setenv("KUBECONFIG", tmp)
	}

	testTable := []struct {
		name     string
		fakePath string
	}{
		{name: "test0", fakePath: "~/.kube/fakeconfig"},
		{name: "test1", fakePath: "~/.notkube/mayberealconfig"},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("KUBECONFIG", tc.fakePath)
			path := getKubeConfig()
			assert.Equal(t, tc.fakePath, path, "they should be equal")
		})
	}
}

// TestGetKubeconfigEnvVar tests getKubeconfig when
// KUBECONFIG environment variable is not set
func TestKubeConfig(t *testing.T) {
	// dont overwrite the kubeconfig path if it exists
	var tmp string
	if os.Getenv("KUBECONFIG") != "" {
		os.Unsetenv("KUBECONFIG")
		// re-write upon exit
		defer os.Setenv("KUBECONFIG", tmp)
	}

	expectedPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	actualPath := getKubeConfig()
	assert.Equal(t, expectedPath, actualPath, "they should be equal")
}
