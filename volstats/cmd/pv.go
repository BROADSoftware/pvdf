package cmd

import (
	"fmt"
	"github.com/BROADSoftware/pvdf/shared/pkg/clientgo"
	"github.com/BROADSoftware/pvdf/shared/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var log = logging.Log.WithFields(logrus.Fields{})

var pvCmd = &cobra.Command{
	Use: "pv",
	Short: "List persistentVolumes and associated usage",
	Run: func(cmd *cobra.Command, args[]string){
		clientSet := clientgo.GetClientSet()
		pvList, err := clientSet.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
		if err != nil {
			panic(fmt.Sprintf("Unable to fetch PersistentVolume list: %v\n", err))
		}
		for _, pv := range pvList.Items {
			log.Infof("PV:%s", pv.Name)
		}


	},

}

func init() {
	rootCmd.AddCommand(pvCmd)
}
