package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/BROADSoftware/pvdf/shared/pkg/clientgo"
	"github.com/BROADSoftware/pvdf/volstats/pkg/lib"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"text/tabwriter"
)

var tplvmCmd = &cobra.Command{
	Use: "tplvm",
	Short: "List Topolvm deviceClass per node",
	Run: func(cmd *cobra.Command, args[]string){
		clientSet := clientgo.GetClientSet()
		tpLvmList := lib.NewTpLvmList(clientSet)
		if len(tpLvmList) > 0 {
			sortTpLvm(tpLvmList)
			if format == "text" {
				tw := new(tabwriter.Writer)
				tw.Init(os.Stdout, 8, 8, 1, '\t', 0)
				_, _ = fmt.Fprintf(tw, "STORAGE CLASS\tDEVICE CLASS\tFSTYPE\tNODE\tSIZE\tFREE\t%%USED")
				for _, tpLvm := range tpLvmList {
					_, _ = fmt.Fprintf(tw, "\n%s\t%s\t%s\t%s\t%s\t%s\t%s", tpLvm.StorageClass, tpLvm.DeviceClass, tpLvm.Fstype, tpLvm.Node, bytes2human(tpLvm.Size, unit), bytes2human(tpLvm.Free, unit), percentToString(tpLvm.Used_pc))
				}
				_, _ = fmt.Fprintf(tw, "\n")
				_ = tw.Flush()
			} else if format == "json" {
				js, err := json.Marshal(tpLvmList)
				if err != nil {
					log.Errorf("Unable to marshal result to json!!")
				} else {
					fmt.Print(string(js))
				}

			} else {
				fmt.Printf("Unknow format ??")
			}
		} else {
			fmt.Printf("No Topovlm storage class !\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(tplvmCmd)
}

func sortTpLvm(l lib.TpLvmList) {
	sort.Slice(l, func(i, j int) bool {
		// Sort order is storageClass/node
		if l[i].StorageClass != l[j].StorageClass {
			return l[i].StorageClass < l[j].StorageClass
		} else {
			return l[i].Node < l[j].Node
		}
	})
}

