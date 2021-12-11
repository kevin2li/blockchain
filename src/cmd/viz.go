package cmd

import (
	"blockchain/src/pkg"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func VizSetup() *cobra.Command {
	var (
		dataset_path string // all_txs.json
		addr_path    string // 1KFHE7w8BhaENAswwryaoccDb6qcT6DbYY
	)
	var vizCmd = &cobra.Command{
		Use:   "viz -d [dataset_path] -a [address_path]",
		Short: "viz entity'addresses in given transcation dataset",
		Long:  `viz entity'addresses(with space splited) in given transcation dataset.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("you don't need to give any argument!")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			t1 := time.Now()
			log.Println("Started!")
			// pkg.Visualize()
			fmt.Printf("dataset path: %+v\n", dataset_path)
			fmt.Printf("address path %+v\n", addr_path)
			err := pkg.StartViz(dataset_path, addr_path)
			if err != nil {
				log.Fatal(err)
			}
			t2 := time.Now()
			log.Println("Finished!")
			fmt.Printf("Time elapsed: %.2f minutes\n", t2.Sub(t1).Minutes())
		},
	}
	vizCmd.Flags().StringVarP(&dataset_path, "dataset_path", "d", "", "path to load transcation dataset")
	vizCmd.Flags().StringVarP(&addr_path, "address_path", "a", "", "path to load entity's address")
	vizCmd.MarkFlagRequired("dataset_path")
	vizCmd.MarkFlagRequired("address_path")
	return vizCmd
}
