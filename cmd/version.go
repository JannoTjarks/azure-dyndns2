package cmd

import (
	"fmt"
	"github.com/JannoTjarks/azure-dyndns2/internal/utils"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  "Show the current version of azure-dyndns2",
	Run: func(cmd *cobra.Command, args []string) {
		if json {
			fmt.Println(utils.GenerateVersionJson())
			return
		}
		
		fmt.Println(utils.GenerateVersionSignature())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&json, "json", false, "If set, the output will be formatted as json")
}
