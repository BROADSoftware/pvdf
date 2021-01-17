package main

import (
	"flag"
	"fmt"
	"github.com/BROADSoftware/pvdf/pkg/logging"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)


var log = logging.Log.WithFields(logrus.Fields{})

func getClientSet(kubeconfig string) *kubernetes.Clientset {
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

func main() {
	logLevel := flag.String("logLevel", "INFO", "Log message verbosity")
	logJson := flag.Bool("logJson", false, "logs in JSON")
	kubeconfig := flag.String("kubeconfig", "", "kubeconfig file")
	flag.Parse()

	logging.ConfigLogger(*logLevel, *logJson)
	log.Info("pvdf start")

	clientSet := getClientSet(*kubeconfig)
	pvList, err := clientSet.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		panic(fmt.Sprintf("Unable to fetch PersistentVolume list: %v\n", err))
	}
	for _, pv := range pvList.Items {
		fmt.Printf("PV:%s\n", pv.Name)
	}
}

