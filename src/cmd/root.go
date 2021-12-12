package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func RootSetup() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "tool",
		Short: "BTC data analysis tool",
		Long: `BTC data analysis tool, including data download, addrress cluster, transaction visualization, etc.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			fmt.Println("Please give subcommond!")
		},
	}
	downloadCmd := CrawlSetup()
	clusterCmd := ClusterSetup() // go run main.go cluster -f data/block_height=711900-711999.json 3Jx1ThGhh5P9vL5XkMw1NauH2YEDSZo4Wd
	vizCmd := VizSetup()
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(vizCmd)
	return rootCmd
}

func Execute() {
	rootCmd := RootSetup()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
