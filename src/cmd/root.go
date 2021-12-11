package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func RootSetup() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "tool",
		Short: "bitcoin data analysis tool",
		Long: `bitcoin data analysis tool, including data download, addrress cluster, transaction visualization, etc.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			fmt.Println("Please give subcommond!")
		},
	}
	downloadCmd := CrawlSetup()
	clusterCmd := ClusterSetup()
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(clusterCmd)
	return rootCmd
}

func Execute() {
	rootCmd := RootSetup()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
