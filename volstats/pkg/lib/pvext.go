package lib

import (
	"github.com/BROADSoftware/pvdf/shared/pkg/logging"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var log = logging.Log.WithFields(logrus.Fields{})

type PvExt struct {
	pv        *v1.PersistentVolume
	pvc       *v1.PersistentVolumeClaim
	pod       *v1.Pod
	Name 	  string
	Namespace string
	Node      string
	Capacity  string
	PodName   string
}

type PvExtList []PvExt

func NewPvExtList(clientSet *kubernetes.Clientset) PvExtList {
	pvList, err := clientSet.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Unable to fetch PersistentVolume list: %v\n", err)
	}
	pvExtList := make(PvExtList, len(pvList.Items))
	// Warning: Don't do _, pv := range.... as pv: &pv will be overrided
	for i, _ := range pvList.Items {
		pvExtList[i] = PvExt{
			pv: &pvList.Items[i],
		}
		//pvExtList[i].Name = pvExtList[i].pv.Name
		pvExtList[i].fillinName()
		pvExtList[i].fillinNode()
		pvExtList[i].fillinCapacity()
	}
	pvExtList.fillinPvc(clientSet)
	pvExtList.fillinPod(clientSet)
	return pvExtList
}
func (this *PvExt) fillinName() {
	this.Name = this.pv.Name
}

func (this *PvExt) fillinNode() {
	if this.pv.Spec.NodeAffinity != nil &&
		this.pv.Spec.NodeAffinity.Required != nil &&
		len(this.pv.Spec.NodeAffinity.Required.NodeSelectorTerms) > 0 &&
		len(this.pv.Spec.NodeAffinity.Required.NodeSelectorTerms[0].MatchExpressions) > 0 &&
		len(this.pv.Spec.NodeAffinity.Required.NodeSelectorTerms[0].MatchExpressions[0].Values) > 0 {
		this.Node = this.pv.Spec.NodeAffinity.Required.NodeSelectorTerms[0].MatchExpressions[0].Values[0]
	}
}


func (this *PvExt) fillinCapacity() {
	strg, ok := this.pv.Spec.Capacity["storage"]
	if ok {
		this.Capacity = (&strg).String()
	}
}

func (this PvExtList) fillinPvc(clientSet *kubernetes.Clientset) {
	pvcList, err := clientSet.CoreV1().PersistentVolumeClaims("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Unable to fetch PersistentVolume list: %v\n", err)
	}
	claimByNsName := make(map[string]v1.PersistentVolumeClaim)
	for _, pvc := range pvcList.Items {
		claimByNsName[ pvc.Namespace + "/" + pvc.Name] = pvc
	}
	for i, _ := range this {
		//log.Debugf("PV: %s   claimRef:%s", this[i].pv.Name, this[i].pv.Spec.ClaimRef.Name)
		if this[i].pv.Spec.ClaimRef != nil {
			this[i].Namespace = this[i].pv.Spec.ClaimRef.Namespace
			pvc, ok := claimByNsName[this[i].Namespace + "/" + this[i].pv.Spec.ClaimRef.Name]
			if ok {
				//log.Debugf("Set pvc %s", pvc.Name)
				this[i].pvc = &pvc
			} else {
				log.Debugf("Orphean pv: '%s", this[i].pv.Name)
			}
		}
	}
}

func (this PvExtList) fillinPod(clientSet *kubernetes.Clientset) {
	pvExtByClaim := make(map[string]*PvExt)
	namespaces := make(map[string]bool)
	for i, _ := range this {
		if this[i].pvc != nil && this[i].Namespace != "" {
			namespaces[this[i].Namespace] = true
			pvExtByClaim[this[i].Namespace + "/" + this[i].pvc.Name] = &this[i]
		}
	}
	for namespace, _ := range namespaces {
		log.Debugf("Load pods from namespace '%s'", namespace)
		podList, err := clientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Unable to fetch pod list in namespace '%s': %v\n", namespace, err)
		}
		for i, pod := range podList.Items {
			for _, volume := range pod.Spec.Volumes {
				if volume.PersistentVolumeClaim != nil {
					claimName := namespace + "/" + volume.PersistentVolumeClaim.ClaimName
					pvExt, ok := pvExtByClaim[claimName]
					if ok {
						pvExt.pod = &podList.Items[i]
						pvExt.PodName = pod.Name
					} else {
						log.Errorf("Unable to find a PV matching claim name '%s' in namespace '%s'", claimName, namespace)
					}
				}
			}

		}
	}
}

