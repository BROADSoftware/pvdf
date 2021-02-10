package main

import (
	"flag"
	"fmt"
	"github.com/BROADSoftware/pvdf/pvscanner/pkg/lib"
	"github.com/BROADSoftware/pvdf/pvscanner/pkg/topolvm"
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
	var version bool
	var nodeName string
	var isTopolvm bool
	var lvmdConfigPath string
	var statfsTimeout string
	var period string
	flag.BoolVar(&version, "version", false, "Display current version")
	flag.StringVar(&logging.Level, "logLevel", "INFO", "Log message verbosity")
	flag.BoolVar(&logging.LogJson, "logJson", false, "logs in JSON")
	flag.StringVar(&clientgo.Kubeconfig, "kubeconfig", "", "kubeconfig file")
	flag.StringVar(&lib.RootfsPath, "rootFsPath", "/", "root FS path")
	flag.StringVar(&statfsTimeout, "statFsTimeout", "5s", "Timeout on syscall failure")
	flag.StringVar(&period, "period", "60s", "Scan period")
	flag.BoolVar(&isTopolvm, "topolvm", false, "Lookup Topolvm LVM VolumeGroup size")
	flag.StringVar(&nodeName, "nodeName", "", "Node name")
	flag.StringVar(&lvmdConfigPath, "lvmdConfigPath", "/etc/topolvm/lvmd.yaml", "Topolvm/lvmd config file path")
	flag.BoolVar(&topolvm.Containerized, "containerized", false, "Is running inside container")
	flag.StringVar(&topolvm.Nsenter, "nsenter", "/usr/bin/nsenter", "nsenter unix command path")
	flag.StringVar(&topolvm.Lvm, "lvm", "/sbin/lvm", "lvm unix command path")
	flag.Parse()

	logging.ConfigLogger()

	if version {
		fmt.Printf("Version: %s\n", lib.Version)
		os.Exit(0)
	}

	var err error
	if lib.StatfsTimeout, err = time.ParseDuration(statfsTimeout); err != nil {
		log.Fatalf("Value '%s' is invalid as a duration for statFsTimeout paramter", statfsTimeout)
	}
	if lib.Period, err = time.ParseDuration(period); err != nil {
		log.Fatalf("Value '%s' is invalid as a duration for period paramter", period)
	}
	if isTopolvm {
		if nn := os.Getenv("NODE_NAME"); nn != "" {
			nodeName = nn
		}
		if nodeName == "" {
			log.Fatalf("NODE_NAME env variable is not defined, no --nodname parameter and --vgsd is set")
		}
	}
	log.Infof("pvscanner start. version:%s. logLevel:%s. Will scan PV every %s", lib.Version, logging.Level, period)

	clientSet := clientgo.GetClientSet()
	var lvmdConfig *topolvm.LvmdConfig
	if isTopolvm {
		lvmdConfig, err = topolvm.LoadLvmdConfig(lvmdConfigPath)
		if err != nil {
			log.Warnf("Unable to load Topolvm lvmd config file '%s':%v. May be Topolvm is not deployed on this node.", lvmdConfigPath, err)
			isTopolvm = false
		}
	}
	for true {
		log.Debugf("-------------------------------------------------")
		workOnPv(clientSet)
		if isTopolvm {
			workOnNode(clientSet, nodeName, lvmdConfig)
		}
		time.Sleep(lib.Period)
	}
}

func adjustAnnotation(node *v1.Node, annotation string, value string) bool {
	if v, ok := node.Annotations[annotation]; (!ok || v != value) {
		log.Infof("Udpate usage information %s:%s", annotation, value)
		node.Annotations[annotation] = string(value)
		return true
	} else {
		log.Debugf("Annotation %s untouched", annotation)
		return false
	}
}

func removeTraingB(x string) string {
	if strings.HasSuffix(x, "B") {
		return x[:len(x)-1]
	} else {
		return x
	}
}

func workOnNode(clientset *kubernetes.Clientset, nodeName string, lvmdConfig *topolvm.LvmdConfig) {
	node, err := clientset.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Unable to load node '%s':%v. VolumeGroup size information will not be updated", nodeName, err)
		return
	}
	vgByName, err := topolvm.GetVgByName()
	if err != nil {
		log.Warnf("Unable to read VolumeGroup size:%v. VolumeGroup size will be unknown", err)
		// Will cleanup all related annotation
		dirty := false
		for k, _ := range node.Annotations {
			if strings.HasPrefix(k, common.SizeTopolvmAnnotationPrefix) {
				delete(node.Annotations, k)
				dirty = true
			}
		}
		if dirty {
			log.Infof("Cleanup all %s/* annotations", common.SizeTopolvmAnnotationPrefix )
			_, err = clientset.CoreV1().Nodes().Update(node)
			if err != nil {
				log.Errorf("Unable to udpate usage information on Node '%s': %v", node.Name, err)
			}
		}
		return
	}
	dirty := false
	dccount := 0
	for _, dc := range lvmdConfig.DeviceClasses {
		vg, ok := vgByName[dc.VolumeGroup]
		if ok {
			dirty = adjustAnnotation(node, fmt.Sprintf(common.SizeTopolvmAnnotation, dc.Name), removeTraingB(vg.VgSize) ) || dirty	// Warning, order is important.
			dccount++
		} else {
			log.Warnf("Unable to find volumeGroup '%s' for deviceClass '%s'", dc.VolumeGroup, dc.Name)
		}
	}
	if dirty {
		_, err = clientset.CoreV1().Nodes().Update(node)
		if err != nil {
			log.Errorf("Unable to udpate usage information on Node '%s': %v", node.Name, err)
		}
	}
	log.Infof("%d Topolvm deviceClasses has been found", dccount)
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

