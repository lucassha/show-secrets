package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestCreateDecodedObject tests struct creation of decoded kubernetes secrets
// into a DecodedSecretList object
func TestCreateDecodedObject(t *testing.T) {
	testTable := []struct {
		name    string
		secrets []apiv1.Secret
		decoded DecodedSecretList
	}{
		{
			name: "test case 1",
			secrets: []apiv1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
					Data: map[string][]byte{
						"data": []byte("thisisnotrealdata"),
					},
				},
			},
			decoded: DecodedSecretList{
				Items: 1,
				DecodedSecrets: []DecodedSecret{
					{
						Name: "test1",
						Data: map[string]string{
							"data": "thisisnotrealdata",
						},
					},
				},
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			got := createDecodedObject(tc.secrets)
			assert.Equal(t, tc.decoded, got, "they should be equal")
		})
	}

}
