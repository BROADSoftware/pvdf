package cmd

import (
	"github.com/BROADSoftware/pvdf/shared/common"
	"github.com/BROADSoftware/pvdf/shared/pkg/clientgo"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var cleanCmd = &cobra.Command{
	Use: "cleanup",
	Short: "Cleanup all annotations set by pvscanner",
	Hidden: false,
	Run: func(cmd *cobra.Command, args[]string){
		clientSet := clientgo.GetClientSet()
		pvList, err := clientSet.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Unable to fetch PersistentVolume list: %v\n", err)
		}
		for _, pv := range pvList.Items {
			dirty := false
			for k, _ := range pv.Annotations {
				if strings.Contains(k, common.RootAnnotation) {
					delete(pv.Annotations, k)
					dirty = true
				}
			}
			if dirty {
				_, err = clientSet.CoreV1().PersistentVolumes().Update(&pv)
				if err != nil {
					log.Errorf("Unable to cleanup usage information on PV '%s':%v", pv.Name, err)
				} else {
					log.Infof("Removed usage information for PV '%s'", pv.Name)
				}
			}
		}
		nodeList, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Unable to fetch Node list: %v\n", err)
		}
		for _, node := range nodeList.Items {
			dirty := false
			for k, _ := range node.Annotations {
				if strings.Contains(k, common.RootAnnotation) {
					delete(node.Annotations, k)
					dirty = true
				}
			}
			if dirty {
				_, err = clientSet.CoreV1().Nodes().Update(&node)
				if err != nil {
					log.Errorf("Unable to cleanup usage information on Node '%s': %v", node.Name, err)
				} else {
					log.Infof("Removed all pvdf related annotations on node '%s'", node.Name)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

