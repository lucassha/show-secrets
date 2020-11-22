package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type flagOptions struct {
	namespace  *string
	kubeconfig *string
}

var flags flagOptions

func main() {
	flags.namespace = flag.String("namespace", apiv1.NamespaceDefault, "namespace to search for secrets")
	flags.kubeconfig = flag.String("kubeconfig", getKubeConfig(), "path to search for kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *flags.kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Could not create client from kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Could not create clientset from config: %v", err)
	}

	// this will replace the lines below
	list, err := getSecrets(clientset, *flags.namespace)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	fmt.Fprintln(w, "NAME\tKEY\tVALUE")

	for _, s := range list.Items {
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

// getSecrets creates the secrets clients from corev1 api
// and returns all secrets in the designated namespace
// func getSecrets(c *kubernetes.Clientset, namespace string) (*apiv1.SecretList, error) {
func getSecrets(c kubernetes.Interface, namespace string) (*apiv1.SecretList, error) {
	secretsClient := c.CoreV1().Secrets(namespace)
	list, err := secretsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return &apiv1.SecretList{}, fmt.Errorf("Could not list secrets: %v", err)
	}
	return list, nil
}

// getKubeConfig passes back the user's kubeconfig file and will
// only be called when --kubeconfig is not overridden
// first precendence is KUBECONFIG env var
// if not set, return $HOME/.kube/config
func getKubeConfig() string {
	kCfgPath := os.Getenv("KUBECONFIG")
	if kCfgPath != "" {
		return kCfgPath
	}
	return filepath.Join(homedir.HomeDir(), ".kube", "config")
}
