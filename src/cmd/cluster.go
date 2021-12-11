package cmd

import (
	"blockchain/src/pkg"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func ClusterSetup() *cobra.Command {
	var (
		dataset_path string // all_txs.json
		start_addr   string // 1KFHE7w8BhaENAswwryaoccDb6qcT6DbYY
	)
	var clusterCmd = &cobra.Command{
		Use:   "cluster -f [dataset_path] [address]",
		Short: "cluster address in given transcation dataset",
		Long:  `cluster address in given transcation dataset.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("you should only give one argument!")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			t1 := time.Now()
			log.Println("Started!")
			start_addr = args[0]
			pkg.StartCluster(dataset_path, start_addr)
			t2 := time.Now()
			log.Println("Finished!")
			fmt.Printf("Time elapsed: %.2f minutes\n", t2.Sub(t1).Minutes())
		},
	}
	clusterCmd.Flags().StringVarP(&dataset_path, "dataset_path", "f", "", "path to load transcation dataset")
	clusterCmd.MarkFlagRequired("dataset_path")
	return clusterCmd
}
