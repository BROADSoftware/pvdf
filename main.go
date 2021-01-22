package main

import (
	"flag"
	"fmt"
	"github.com/BROADSoftware/pvdf/pkg/lib"
	"github.com/BROADSoftware/pvdf/pkg/logging"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)


var log = logging.Log.WithFields(logrus.Fields{})


func main() {
	logLevel := flag.String("logLevel", "DEBUG", "Log message verbosity")
	logJson := flag.Bool("logJson", false, "logs in JSON")
	kubeconfig := flag.String("kubeconfig", "", "kubeconfig file")
	flag.Parse()

	logging.ConfigLogger(*logLevel, *logJson)
	log.Info("pvdf start")

	clientSet := lib.GetClientSet(*kubeconfig)
	for true {
		log.Infof("------------------------")
		work(clientSet)
		time.Sleep(lib.Period)
	}
}


func work(clientSet *kubernetes.Clientset) {
	// Retrieve all PV from api server
	pvList, err := clientSet.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		panic(fmt.Sprintf("Unable to fetch PersistentVolume list: %v\n", err))
	}

	// Get all mounted file system on this node
	fileSystems, err := lib.ListFileSystems()
	if err != nil {
		panic(fmt.Sprintf("Unable to fetch Mountpoints: %v\n", err))
	}
	// Filter mounted fs and take only the ones which are pod volume.
	volumeByName := lib.GetVolumeByName(fileSystems)

	factor := uint64(1024*1024) // Mb
	// And now, loop for all PV to find matching volumes and populate information
	for _, pv := range pvList.Items {
		//fmt.Printf("PV:%s\n", pv.Name)
		volume, ok := volumeByName[pv.Name]
		if ok {
			volume.GetStats()
			if volume.Stats.Error == nil {
				log.Debugf("PV '%s':  size:%d   free:%d  avail:%d", volume.Name, volume.Stats.Size/factor, volume.Stats.Free/factor, volume.Stats.Avail/factor)
			} else {
				log.Warnf("PV: '%s': Error:%s  ", volume.Name, volume.Stats.Error)
			}
		} else {
			log.Tracef("No volume for PV '%s'. Should be on another node", pv.Name)
		}
	}
}