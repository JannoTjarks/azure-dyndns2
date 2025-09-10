package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	hostname          string
	myip              string
	dnsZoneName       string
	resourceGroupName string
	subscriptionId    string
	port              string
	json              bool
)

var rootCmd = &cobra.Command{
	Use:   "azure-dyndns2",
	Short: "Simple dyndns2-compatible web api for Azure DNS",
	Long: `Simple dyndns2-compatible web api for Azure DNS
---
This project goal is to enable Azure DNS to support DynDNS by 
adding a small webapi, which can be addressed by the DynDNS 
Update API, also called the dyndns2 standard. That’s why this 
project is called azure-dyndns2.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Moin World")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
