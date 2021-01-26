package cmd

import (
	"fmt"
	"github.com/BROADSoftware/pvdf/shared/pkg/clientgo"
	"github.com/BROADSoftware/pvdf/volstats/pkg/lib"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)



var pvCmd = &cobra.Command{
	Use: "pv",
	Short: "List persistentVolumes and associated usage",
	Run: func(cmd *cobra.Command, args[]string){
		clientSet := clientgo.GetClientSet()
		pvExtList := lib.NewPvExtList(clientSet)
		if len(pvExtList) > 0 {
			tw := new(tabwriter.Writer)
			tw.Init(os.Stdout, 8, 8, 1, '\t', 0)
			_, _ = fmt.Fprintf(tw, "NAMESPACE\tNODE\tPV NAME\tPOD NAME\tCAP.")
			for _, pvExt := range pvExtList {
				_, _ = fmt.Fprintf(tw, "\n%s\t%s\t%s\t%s\t%s", pvExt.Namespace, pvExt.Node, pvExt.Name, pvExt.PodName, pvExt.Capacity)
			}
			_, _ = fmt.Fprintf(tw, "\n")
			_ = tw.Flush()
		} else {
			fmt.Printf("No PersistentVolume !\n")
		}
	},

}

func init() {
	rootCmd.AddCommand(pvCmd)
}
