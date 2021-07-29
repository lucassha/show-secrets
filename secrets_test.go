package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func CreateFakeClientAndData(t testing.TB) *testclient.Clientset {
	t.Helper()
	clientset := testclient.NewSimpleClientset()

	testData := []struct {
		secret *apiv1.Secret
		data   []byte
		name   string
	}{
		{data: []byte("definitelynotfakedata"), name: "fakesecret0"},
		{data: []byte("supersecret"), name: "fakesecret1"},
	}

	for i := 0; i < 2; i++ {
		testData[i].secret = &apiv1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testData[i].name,
				Namespace: apiv1.NamespaceDefault,
			},
			Data: map[string][]byte{
				"secret": []byte(testData[i].data),
			},
		}
	}

	for _, s := range testData {
		// Add secrets to fake client
		_, err := clientset.CoreV1().Secrets(apiv1.NamespaceDefault).Create(context.TODO(), s.secret, metav1.CreateOptions{})
		if err != nil {
			t.Fatal(err)
		}
	}

	return clientset
}

// TestGetSecrets creates a fake client and mock list data
// in order to confirm all secrets are returned
func TestGetSecrets(t *testing.T) {
	// create fake client and two secrets in helper function
	client := CreateFakeClientAndData(t)
	list, err := getSecrets(client, apiv1.NamespaceDefault)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name       string
		key        string
		secretData []byte
	}{
		{"test0", "secret", []byte("definitelynotfakedata")},
		{"test1", "secret", []byte("supersecret")},
	}

	// check each secret exists and eval the two []bytes
	for i, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, list.Items[i].Data[tc.key], tc.secretData, "should be equal")
		})
	}
}
