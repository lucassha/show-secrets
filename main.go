package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type flagOptions struct {
	namespace  *string
	kubeconfig *string
	o          *string
	output     *string
}

const (
	defaultOutput   = "wide"
	jsonOutput      = "json"
	outputUsage     = "output format"
	namespaceUsage  = "namespace to search for secret(s)"
	kubeconfigUsage = "path to search for kubeconfig file"
)

var flags flagOptions

func init() {
	flags.namespace = flag.String("namespace", apiv1.NamespaceDefault, namespaceUsage)
	flags.namespace = flag.String("n", apiv1.NamespaceDefault, namespaceUsage+" (shorthand)")
	flags.kubeconfig = flag.String("kubeconfig", getKubeConfig(), kubeconfigUsage)
	flags.output = flag.String("output", defaultOutput, outputUsage)
	flags.o = flag.String("o", defaultOutput, outputUsage+" (shorthand)")
}

func main() {
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *flags.kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Could not create client from kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Could not create clientset from config: %v", err)
	}

	list, err := getSecrets(clientset, *flags.namespace)

	// print output based on -o or --output flag
	switch {
	case *flags.o == jsonOutput:
		jsonPrintSecrets(list.Items)
	default:
		widePrintSecrets(list.Items)
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
