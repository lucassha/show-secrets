package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	apiv1 "k8s.io/api/core/v1"
)

type Secret struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DecodedSecret struct {
	Name    string   `json:"name"`
	Secrets []Secret `json:"secrets"`
}

type DecodedSecretList struct {
	Items          int             `json:"items"`
	DecodedSecrets []DecodedSecret `json:"decodedSecrets"`
}

func createJsonObject(secrets []apiv1.Secret) DecodedSecretList {
	d := DecodedSecretList{}

	// instantiate a new block of memory for slice of decodedSecrets
	d.DecodedSecrets = make([]DecodedSecret, len(secrets))
	d.Items = len(secrets)

	for i, s := range secrets {
		d.DecodedSecrets[i].Name = s.GetName()
		j := 0

		// instantiate a new block of memory for slice of secrets
		d.DecodedSecrets[i].Secrets = make([]Secret, len(s.Data))

		for k, v := range s.Data {
			d.DecodedSecrets[i].Secrets[j].Key = string(k)
			d.DecodedSecrets[i].Secrets[j].Value = string(v)
			j++
		}
	}

	return d
}

// jsonPrintSecrets takes a list of secrets and prints the
// decoded secret output in json, similar to `kubectl get * -o json`
func jsonPrintSecrets(secrets []apiv1.Secret) {
	d := createJsonObject(secrets)
	b, _ := json.MarshalIndent(&d, "", "  ")
	fmt.Println(string(b))
}

// widePrintSecrets takes a list of secrets and prints the
// decoded secret output in a wide output, similar to `kubectl get *`
func widePrintSecrets(secrets []apiv1.Secret) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	fmt.Fprintln(w, "NAME\tKEY\tVALUE")

	for _, s := range secrets {
		// dont show the default service account token secret
		if strings.Contains(s.ObjectMeta.GetName(), "default-token") {
			continue
		}

		i := 1
		secretName := s.GetName()

		for k, v := range s.Data {
			switch {
			case i == 1:
				fmt.Fprintf(w, "%s\t%s\t%s\n", secretName, k, string(v))
			case i == len(s.Data):
				fmt.Fprintf(w, "└── %s\t%s\t%s\n", "", k, string(v))
			default:
				fmt.Fprintf(w, "├── %s\t%s\t%s\n", "", k, string(v))
			}

			i++
		}
		w.Flush()
	}
}
