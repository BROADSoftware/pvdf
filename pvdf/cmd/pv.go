package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/BROADSoftware/pvdf/pvdf/pkg/lib"
	"github.com/BROADSoftware/pvdf/shared/pkg/clientgo"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

var noScannerCheck bool

var pvCmd = &cobra.Command{
	Use:   "pv",
	Short: "List persistentVolumes and associated usage",
	Run: func(cmd *cobra.Command, args []string) {
		clientSet := clientgo.GetClientSet()
		pvExtList := lib.NewPvExtList(clientSet)
		if len(pvExtList) > 0 {
			if !noScannerCheck {
				if !checkDaemon(clientSet, pvExtList) {
					log.Fatalf("'pvscanner' daemonset seemns not being running properly. Data will be invalid. You can skip this test by providing --noScannerCheck option\n")
					os.Exit(2)
				}
			}
			sortPv(pvExtList)
			if format == "text" {
				tw := new(tabwriter.Writer)
				tw.Init(os.Stdout, 8, 8, 1, '\t', 0)
				_, _ = fmt.Fprintf(tw, "NAMESPACE\tNODE\tPV NAME\tPOD NAME\tREQ.\tSTORAGE CLASS\tSIZE\tFREE\t%%USED")
				for _, pvExt := range pvExtList {
					_, _ = fmt.Fprintf(tw, "\n%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s", pvExt.Namespace, pvExt.Node, shorten(pvExt.Name, 20), shorten(pvExt.PodName, 20), pvExt.Capacity, pvExt.StorageClass, bytes2human(pvExt.Size, unit), bytes2human(pvExt.Free, unit), percentToString(pvExt.Used_pc))
				}
				_, _ = fmt.Fprintf(tw, "\n")
				_ = tw.Flush()
			} else if format == "json" {
				js, err := json.Marshal(pvExtList)
				if err != nil {
					log.Errorf("Unable to marshal result to json!!")
				} else {
					fmt.Print(string(js))
				}
			} else {
				fmt.Printf("Unknow format ??")
			}
		} else {
			fmt.Printf("No PersistentVolume !\n")
		}
	},
}

func init() {
	pvCmd.PersistentFlags().BoolVarP(&noScannerCheck, "noScannerCheck", "n", false, "Don't perform security pvscanner daemonset check." )
	rootCmd.AddCommand(pvCmd)

}

func sortPv(l lib.PvExtList) {
	sort.Slice(l, func(i, j int) bool {
		// Sort order is namespace/storageclass/pod/node/name
		if l[i].Namespace != l[j].Namespace {
			return l[i].Namespace < l[j].Namespace
		} else {
			if l[i].StorageClass != l[j].StorageClass {
				return l[i].StorageClass < l[j].StorageClass
			} else {
				if l[i].PodName != l[j].PodName {
					return l[i].PodName < l[j].PodName
				} else {
					if l[i].Node != l[j].Node {
						return l[i].Node < l[j].Node
					} else {
						return l[i].Name < l[j].Name
					}
				}
			}
		}
	})
}

// This is a security. We check a pvscanner pod in the pvdf-system exists on each concerned node.
// If not there is big chance stats will be obsoleted.
func checkDaemon(clientSet *kubernetes.Clientset, pvExtList lib.PvExtList) bool {
	podList, err := clientSet.CoreV1().Pods("pvdf-system").List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("Unable to fetch pods in pvdf-system namespace: '%v'", err)
		return false
	}
	nodeWithRunningPod := make(map[string]struct{})
	for _, pod := range podList.Items {
		if strings.Contains(pod.Name, "pvscanner") && pod.Status.ContainerStatuses[0].State.Running != nil {
			node := pod.Spec.NodeName
			nodeWithRunningPod[node] = struct{}{}
		}
	}
	nodeSet := pvExtList.GetNodesSet(clientSet)
	for node, _ := range nodeSet {
		if _, ok := nodeWithRunningPod[node]; !ok {
			log.Warnf("Unable to find a 'pvscanner...' pod running on node '%s'", node)
			return false
		}
	}
	return true
}
