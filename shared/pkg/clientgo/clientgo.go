package clientgo

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// To be set by CLI parameters before calling GetClientSet()
var Kubeconfig string

var clientset *kubernetes.Clientset
var once sync.Once

func GetClientSet() *kubernetes.Clientset {
	once.Do(func() {
		var config *rest.Config = nil
		var err error
		if Kubeconfig == "" {
			if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
				Kubeconfig = envvar
			}
		}
		if Kubeconfig == "" {
			config, err = rest.InClusterConfig()
		}
		if config == nil {
			home := homedir.HomeDir()
			if Kubeconfig == "" && home != "" {
				Kubeconfig = filepath.Join(home, ".kube", "config")
			}
			config, err = clientcmd.BuildConfigFromFlags("", Kubeconfig)
			if err != nil {
				log.Fatalf("The kubeconfig cannot be loaded: %v\n", err)
			}
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Unable to instanciate clientset: %v\n", err)
		}
	})
	return clientset
}

