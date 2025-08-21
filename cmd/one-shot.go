package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var oneshotCmd = &cobra.Command{
	Use:     "one-shot",
	Aliases: []string{"oneshot"},
	Short:   "Set a DNS Record direct from the CLI",
	Long:    "Allows you to set a DNS Record direct from the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Demo")
	},
}

func init() {
	rootCmd.AddCommand(oneshotCmd)
}
