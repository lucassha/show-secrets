package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	apiv1 "k8s.io/api/core/v1"
)

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
