package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type AuthInfo struct {
	Host          string `json:"host"`
	CaCertificate string `json:"cluster_ca_certificate"`
	Token         string `json:"token"`
}

type Query struct {
	Namespace     string `json:"namespace"`
	Name          string `json:"name"`
	Key           string `json:"key"`
	Context       string `json:"context"`
	Host          string `json:"host"`
	CaCertificate string `json:"cluster_ca_certificate"`
	Token         string `json:"token"`
}

type Result struct {
	Value string `json:"value"`
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func main() {
	// Get terraform input
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatal("cannot read stdin input: %v", err)
	}

	var q Query
	if err := json.Unmarshal(input, &q); err != nil {
		fatal("cannot json umarshal input: %v", err)
	}

	if q.Namespace == "" {
		fatal("missing or empty namespace parameter")
	}

	if q.Name == "" {
		fatal("missing or empty name parameter")
	}

	if q.Key == "" {
		fatal("missing or empty key parameter")
	}

	// if q.Context == "" {
	// 	fatal("missing or empty context parameter")
	// }

	if q.Host == "" {
		fatal("missing or empty k8s host")
	}

	if q.CaCertificate == "" {
		fatal("missing or empty k8s ca certificate")
	}

	if q.Token == "" {
		fatal("missing or empty k8s token")
	}

	// Build kubernetes configuration
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			ClusterInfo: api.Cluster{
				Server:                   q.Host,
				CertificateAuthorityData: []byte(q.CaCertificate)},
			AuthInfo: api.AuthInfo{
				Token: q.Token,
			}}).ClientConfig()

	if err != nil {
		fatal("cannot load Kubernetes configuration: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fatal("cannot create Kubernetes configuration: %v", err)
	}

	secret, err := clientset.CoreV1().Secrets(q.Namespace).Get(context.TODO(), q.Name, metav1.GetOptions{})
	if err != nil {
		fatal("cannot get the %q secret in %q namespace: %v", q.Name, q.Namespace, err)
	}

	val, ok := secret.Data[q.Key]
	if !ok {
		fatal("cannot found the key %q for %q secret in %q namespace", q.Key, q.Name, q.Namespace)
	}

	r := &Result{Value: string(val)}

	o, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fatal("cannot marshal indent result: %v", err)
	}

	os.Stdout.Write(o)
	// return
}
