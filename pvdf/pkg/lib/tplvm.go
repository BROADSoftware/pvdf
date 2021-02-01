package lib

import (
	"fmt"
	"github.com/BROADSoftware/pvdf/shared/common"
	coreV1 "k8s.io/api/core/v1"
	storageV1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"strings"
)

type TpLvm struct {
	node *coreV1.Node
	storageClass *storageV1.StorageClass
	Free int64
	Size int64
	Used_pc int
	DeviceClass string
	StorageClass string
	Node string
	Fstype string
}

func (this *TpLvm) fillin() {
	this.StorageClass = this.storageClass.Name
	this.Node = this.node.Name
	this.Fstype = this.storageClass.Parameters[StorageClassFstypeKey]
}


type TpLvmList []TpLvm

var TopolvmProvisioner = "topolvm.cybozu.com"
var TopolvmDcParameter = "topolvm.cybozu.com/device-class"
var TopolvmCapacityKey = "capacity.topolvm.cybozu.com"
var StorageClassFstypeKey = "csi.storage.k8s.io/fstype"

func NewTpLvmList(clientSet *kubernetes.Clientset) TpLvmList {
	tpvmList := make(TpLvmList, 0, 10)
	scList, err := clientSet.StorageV1().StorageClasses().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Unable to fetch StorageClass list: %v\n", err)
	}
	scByDevice := make(map[string]*storageV1.StorageClass)
	for i, sc := range scList.Items {
		if sc.Provisioner == TopolvmProvisioner {
			dc, ok := sc.Parameters[TopolvmDcParameter]
			if !ok {
				dc = "00default"
			}
			if sc2, ok := scByDevice[dc]; ok {
				log.Warnf("There is more than one StorageClass for deviceClasse '%s'  (%s and %s)", dc, sc.Name, sc2.Name )
			}
			scByDevice[dc] = &scList.Items[i]
		}
	}

	nodeList, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Unable to fetch Node list: %v\n", err)
	}
	for i, node := range nodeList.Items {
		for k, v := range node.Annotations {
			s := strings.Split(k, "/")
			if len(s) == 2 && s[0] == TopolvmCapacityKey {
				dc := s[1]
				storageClass, ok := scByDevice[dc]
				if !ok {
					if dc != "00default" {
						log.Warnf("Not StorageClass for DeviceClass '%s'", dc)
					}
				} else {
					free, _ := strconv.ParseInt(v, 10, 64)
					sizeKey := fmt.Sprintf(common.SizeTopolvmAnnotation, dc)
					sizeStr, ok := node.Annotations[sizeKey]
					var size int64 = -1
					var used_pc = -1
					if ok {
						size, _ = strconv.ParseInt(sizeStr, 10, 54)
					} else {
						log.Warnf("Unable to find '%s' annotation on this node. Is pvscanner daemon set deployed ?")
					}
					if size != -1 && free != -1 {
						used_pc = int(((size-free)*100)/size)
					}
					tplvm := TpLvm{
						node: &nodeList.Items[i],
						storageClass: storageClass,
						Free: free,
						Size : size,
						Used_pc: used_pc,
						DeviceClass: dc,
					}
					tplvm.fillin()
					tpvmList = append(tpvmList, tplvm)
				}
			}
		}
	}
	return tpvmList
}