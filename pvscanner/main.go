package main

import (
	"flag"
	"fmt"
	"github.com/BROADSoftware/pvdf/pvscanner/pkg/lib"
	"github.com/BROADSoftware/pvdf/pvscanner/pkg/logging"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)


var log = logging.Log.WithFields(logrus.Fields{})


func main() {
	logLevel := flag.String("logLevel", "INFO", "Log message verbosity")
	logJson := flag.Bool("logJson", false, "logs in JSON")
	kubeconfig := flag.String("kubeconfig", "", "kubeconfig file")
	flag.StringVar(&lib.ProcPath, "procPath", "/proc", "proc device path")
	flag.StringVar(&lib.RootfsPath, "rootFsPath", "/", "root FS path")
	statfsTimeout := flag.String("statFsTimeout", "5s", "Timeout on syscall failure")
	period := flag.String("period", "60s", "Scan period")
	flag.Parse()

	logging.ConfigLogger(*logLevel, *logJson)

	var err error
	if lib.StatfsTimeout, err = time.ParseDuration(*statfsTimeout); err != nil {
		log.Fatalf("Value '%s' is invalid as a duration for statFsTimeout paramter", *statfsTimeout)
	}
	if lib.Period, err = time.ParseDuration(*period); err != nil {
		log.Fatalf("Value '%s' is invalid as a duration for period paramter", *period)
	}
	log.Infof("pvscanner start. Will scan PV every %s", *period)

	clientSet := lib.GetClientSet(*kubeconfig)
	for true {
		log.Debugf("------------------------")
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
	pvCount := 0
	for _, pv := range pvList.Items {
		//fmt.Printf("PV:%s\n", pv.Name)
		volume, ok := volumeByName[pv.Name]
		if ok {
			volume.GetStats()
			if volume.Stats.Error == nil {
				log.Debugf("PV '%s':  size_mib:%d  free_mib:%d  Used_mib:%d (%d%%)",
					volume.Name,
					volume.Stats.Size/factor,
					volume.Stats.Free/factor,
					(volume.Stats.Size-volume.Stats.Free)/factor,
					(volume.Stats.Size-volume.Stats.Free)*100/volume.Stats.Size)
				if volume.AdjustAnnotationsOn(pv) {
					_, err = clientSet.CoreV1().PersistentVolumes().Update(&pv)
					if err != nil {
						log.Errorf("Unable to udpate usage information on PV '%s'", volume.Name)
					} else {
						log.Infof("Udpate usage information for PV '%s' (size_mib:%s  free_mib:%s)", volume.Name, pv.Annotations[lib.SizeAnnotation], pv.Annotations[lib.FreeAnnotation])
					}
				}
			} else {
				log.Warnf("PV: '%s': Error:%s  (Usage annotation will be removed)", volume.Name, volume.Stats.Error)
				// In such case, better to remove our annotations
				delete(pv.Annotations, lib.FreeAnnotation)
				delete(pv.Annotations, lib.SizeAnnotation)
				_, err = clientSet.CoreV1().PersistentVolumes().Update(&pv)
				if err != nil {
					log.Errorf("Unable to udpate usage information on PV '%s'", volume.Name)
				}
			}
			pvCount++
		} else {
			log.Tracef("No volume for PV '%s'. Should be on another node", pv.Name)
		}
	}
	log.Infof("%d PVs has been scanned", pvCount)
}