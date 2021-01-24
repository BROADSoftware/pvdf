package cmd

import (
	"fmt"
	"github.com/BROADSoftware/pvdf/shared/pkg/clientgo"
	"github.com/BROADSoftware/pvdf/shared/pkg/logging"
	"github.com/spf13/cobra"
	"os"
)


var rootCmd = &cobra.Command{
	Use:   "pvstatus",
	Short: "A PV usage display tool",
}


func init() {
	rootCmd.PersistentFlags().StringVarP(&clientgo.Kubeconfig, "kubeconfig", "k", "", "kubeconfig file" )
	rootCmd.PersistentFlags().StringVarP(&logging.Level, "loglevel", "l", "INFO", "Log level" )
	rootCmd.PersistentFlags().BoolVarP(&logging.LogJson, "logJson", "j", false, "kubeconfig file" )
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		logging.ConfigLogger()
	}
}



var debug = true

func Execute() {
	defer func() {
		if !debug {
			if r := recover(); r != nil {
				fmt.Printf("ERROR:%v\n", r)
				os.Exit(1)
			}
		}
	}()
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(2)
	}
}
