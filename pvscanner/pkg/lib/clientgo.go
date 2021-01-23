package lib

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func GetClientSet(kubeconfig string) *kubernetes.Clientset {
	var config *rest.Config = nil
	var err error
	if kubeconfig == "" {
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
	}
	if kubeconfig == "" {
		config, err = rest.InClusterConfig()
	}
	if config == nil {
		home := homedir.HomeDir()
		if kubeconfig == "" && home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("The kubeconfig cannot be loaded: %v\n", err)
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Unable to instanciate clientset: %v\n", err)
	}
	return clientset
}

