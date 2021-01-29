package main

import (
	"flag"
	"fmt"
	"github.com/BROADSoftware/pvdf/pvscanner/pkg/lib"
	"github.com/BROADSoftware/pvdf/shared/common"
	"github.com/BROADSoftware/pvdf/shared/pkg/clientgo"
	"github.com/BROADSoftware/pvdf/shared/pkg/logging"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
	"time"
)

var log = logging.Log.WithFields(logrus.Fields{})

func main() {
	var nodeName string
	var noVgScan bool
	flag.StringVar(&logging.Level, "loglevel", "INFO", "Log message verbosity")
	flag.BoolVar(&logging.LogJson, "logjson", false, "logs in JSON")
	flag.StringVar(&clientgo.Kubeconfig, "kubeconfig", "", "kubeconfig file")
	flag.StringVar(&lib.ProcPath, "procpath", "/proc", "proc device path")
	flag.StringVar(&lib.RootfsPath, "rootfspath", "/", "root FS path")
	statfsTimeout := flag.String("statfstimeout", "5s", "Timeout on syscall failure")
	period := flag.String("period", "60s", "Scan period")
	flag.StringVar(&nodeName, "nodename", "", "Node name")
	flag.BoolVar(&noVgScan, "novgscan", false, "No LVM VG scan")
	flag.Parse()

	logging.ConfigLogger()

	var err error
	if lib.StatfsTimeout, err = time.ParseDuration(*statfsTimeout); err != nil {
		log.Fatalf("Value '%s' is invalid as a duration for statFsTimeout paramter", *statfsTimeout)
	}
	if lib.Period, err = time.ParseDuration(*period); err != nil {
		log.Fatalf("Value '%s' is invalid as a duration for period paramter", *period)
	}
	if !noVgScan {
		if nn := os.Getenv("NODE_NAME"); nn != "" {
			nodeName = nn
		}
		if nodeName == "" {
			log.Fatalf("NODE_NAME env variable is not defined, no --nodname parameter and --novgscan is not set")
		}
	}

	log.Infof("pvscanner start. Will scan PV every %s. logLevel is '%s'", *period, logging.Level)

	clientSet := clientgo.GetClientSet()
	for true {
		log.Debugf("-------------------------------------------------")
		workOnPv(clientSet)
		if !noVgScan {
			workOnNode(clientSet, nodeName)
		}
		time.Sleep(lib.Period)
	}
}

func adjustAnnotation(node *v1.Node, annotation string, value string) bool {
	if v, ok := node.Annotations[annotation]; (!ok || v != value) {
		log.Debugf("Update annotation %s = %s", annotation, value)
		node.Annotations[annotation] = string(value)
		return true
	} else {
		log.Debugf("Annotation %s untouched", annotation)
		return false
	}
}

func workOnNode(clientset *kubernetes.Clientset, nodeName string) {
	lvmVgs, err := lib.GetLvmVg()
	if err != nil {
		log.Errorf("Unable to scan LVM VG:%v", err)
		return
	}
	node, err := clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Unable to load node '%s':%v", nodeName, err)
		return
	}
	vgcount := 0
	vgList := ""
	vgListSep := ""
	dirty := false
	for i := 0; i < len(lvmVgs.Report); i++ {
		for j := 0; j < len(lvmVgs.Report[i].Vg); j++ {
			log.Debugf("LVM VolumeGroup:%s  size:%s  free:%s\n", lvmVgs.Report[i].Vg[j].VgName, lvmVgs.Report[i].Vg[j].VgSize, lvmVgs.Report[i].Vg[j].VgFree)
			vgname := lvmVgs.Report[i].Vg[j].VgName
			vgList = vgList + vgListSep + vgname
			vgListSep = ","
			vgcount++
			dirty = dirty || adjustAnnotation(node, fmt.Sprintf(common.NodeVgSizeAnnotation, vgname), lvmVgs.Report[i].Vg[j].VgSize )
			dirty = dirty || adjustAnnotation(node, fmt.Sprintf(common.NodeVgFreeAnnotation, vgname), lvmVgs.Report[i].Vg[j].VgFree )
		}
	}
	dirty = dirty || adjustAnnotation(node, common.NodeVgListAnnotation, vgList )
	if dirty {
		_, err = clientset.CoreV1().Nodes().Update(node)
		if err != nil {
			log.Errorf("Unable to udpate usage information on Node '%s': %v", node.Name, err)
		} else {
			for k, v := range node.Annotations {
				if strings.HasPrefix(k, common.RootAnnotation) {
					log.Infof("Udpate usage information for Node '%s':%s:%s", node.Name, k, v)
				}
			}
		}
	}
	log.Infof("%d LVM VGs has been found", vgcount)
}


func workOnPv(clientSet *kubernetes.Clientset) {
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
						log.Errorf("Unable to udpate usage information on PV '%s':%v", volume.Name, err)
					} else {
						log.Infof("Udpate usage information for PV '%s' (size_mib:%s  free_mib:%s)", volume.Name, pv.Annotations[common.PvSizeAnnotation], pv.Annotations[common.PvFreeAnnotation])
					}
				}
			} else {
				log.Warnf("PV: '%s': Error:%s  (Usage annotation will be removed)", volume.Name, volume.Stats.Error)
				// In such case, better to remove our annotations
				delete(pv.Annotations, common.PvFreeAnnotation)
				delete(pv.Annotations, common.PvSizeAnnotation)
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