package cmd

import (
	"fmt"
	"log"

	"github.com/JannoTjarks/azure-dyndns2/internal/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var oneshotCmd = &cobra.Command{
	Use:     "one-shot",
	Aliases: []string{"oneshot"},
	Short:   "Set a DNS Record direct from the CLI",
	Long:    "Allows you to set a DNS Record direct from the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		provState, err := utils.CreateOrUpdateDynDnsRecord(
			viper.GetString(hostname.name),
			viper.GetString(myip.name),
			viper.GetString(dnsZoneName.name),
			viper.GetString(dnsResourceGroupName.name),
			viper.GetString(dnsSubscriptionId.name),
		)
		if err != nil {
			log.Fatal("An error occured...")
		}

		fmt.Print(provState)
	},
}

func init() {
	rootCmd.AddCommand(oneshotCmd)

	oneshotCmd.Flags().StringP(hostname.name, hostname.shorthand, hostname.value, hostname.usage)
	oneshotCmd.Flags().StringP(myip.name, myip.shorthand, myip.value, myip.usage)
	oneshotCmd.Flags().String(dnsZoneName.name, dnsZoneName.value, dnsZoneName.usage)
	oneshotCmd.Flags().String(dnsResourceGroupName.name, dnsResourceGroupName.value, dnsResourceGroupName.usage)
	oneshotCmd.Flags().String(dnsSubscriptionId.name, dnsSubscriptionId.value, dnsSubscriptionId.usage)

	oneshotCmd.PreRun = func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())
		viper.AutomaticEnv()

		viper.BindEnv(dnsZoneName.name, dnsZoneName.envVar)
		viper.BindEnv(dnsResourceGroupName.name, dnsResourceGroupName.envVar)
		viper.BindEnv(dnsSubscriptionId.name, dnsSubscriptionId.envVar)
	}
}
