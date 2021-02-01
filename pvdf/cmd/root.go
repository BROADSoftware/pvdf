package cmd

import (
	"fmt"
	"github.com/BROADSoftware/pvdf/shared/pkg/clientgo"
	"github.com/BROADSoftware/pvdf/shared/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var log = logging.Log.WithFields(logrus.Fields{})

var rootCmd = &cobra.Command{
	Use:   "pvdf",
	Short: "A PV usage display tool",
}

var format string // text or json
var unit string // A (Auto), B, K, Ki, M, Mi, G, Gi, T, Ti, P, Pi

func init() {
	rootCmd.PersistentFlags().StringVarP(&clientgo.Kubeconfig, "kubeconfig", "k", "", "kubeconfig file" )
	rootCmd.PersistentFlags().StringVarP(&logging.Level, "logLevel", "l", "INFO", "Log level" )
	rootCmd.PersistentFlags().BoolVarP(&logging.LogJson, "logJson", "j", false, "Logs in JSON" )
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "Output format (text or json)" )
	rootCmd.PersistentFlags().StringVarP(&unit, "unit", "u", "A", "Unit for storage values display" )
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		logging.ConfigLogger()
		format = strings.ToLower(format)
		if format == "txt" {
			format = "text"
		}
		if format != "text" && format != "json" {
			fmt.Printf("Invalid format value (%s). Must be 'text' or 'json'\n", format)
			os.Exit(2)
		}
		if _, ok := factorByUnit[strings.ToLower(unit)]; !ok && strings.ToLower(unit) != "a" && strings.ToLower(unit) != "h"  {
			fmt.Printf("Invalid unit value ('%s'). Must be one of B,Bi,K,Ki,M,Mi,G,Gi,T,Ti,P,Pi\n", unit)
			os.Exit(2)
		}
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

