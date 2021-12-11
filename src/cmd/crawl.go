package cmd

import (
	"blockchain/src/pkg"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func CrawlSetup() *cobra.Command {
	var crawler = pkg.Crawler{
		GetBlockUrl: "https://api.blockchain.info/haskoin-store/btc/block/heights?heights=%s&notx=false",
		GetTxUrl:    "https://api.blockchain.info/haskoin-store/btc/transactions?txids=%s",
		Page:        10,
	}
	crawler.GetUAPool("/home/likai/code/go_program/go_learn2/data/ua.json")

	var (
		isInterval bool
		filepath   string
	)

	var downloadCmd = &cobra.Command{
		Use:   "download [heights to download]",
		Short: "download transactions in given block heights",
		Long: `download transactions in given block heights.
Please give reasonable block heights.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if !isInterval && filepath == "" {
				if len(args) == 0 {
					return errors.New("not enough arguments")
				}
			}
			if isInterval {
				if len(args) != 2 {
					return errors.New("you should only given 2 args with `-r` flag")
				}
				l, _ := strconv.Atoi(args[0])
				r, _ := strconv.Atoi(args[1])
				if l > r {
					return errors.New("the second argument should greater than or equal to the first argument")
				}
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			t1 := time.Now()
			log.Println("Started!")
			// download txs in given block heights range
			if isInterval {
				low, _ := strconv.Atoi(args[0])
				high, _ := strconv.Atoi(args[1])
				blocks, err := crawler.GetBlocksInRange(low, high)
				if err != nil {
					log.Fatalf("%+v\n", err)
				}
				crawler.DownloadAllBlocks(blocks)
				// read heights from file
			} else if filepath != "" {
				heights, _ := pkg.ReadHeights(filepath)
				blocks, err := crawler.GetBlocks(heights)
				if err != nil {
					log.Fatalf("%+v\n", err)
				}
				crawler.DownloadAllBlocks(blocks)
				// download txs in given heights
			} else {
				heights, _ := pkg.Strings2Ints(args)
				blocks, err := crawler.GetBlocks(heights)
				if err != nil {
					log.Fatalf("%+v\n", err)
				}
				crawler.DownloadAllBlocks(blocks)
			}
			t2 := time.Now()
			log.Println("Finished!")
			fmt.Printf("Time elapsed: %.2f minutes\n", t2.Sub(t1).Minutes())
		},
	}

	downloadCmd.Flags().BoolVarP(&isInterval, "interval", "r", false, "")
	downloadCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "file store heights to download")
	downloadCmd.Flags().StringVarP(&crawler.Savedir, "savedir", "s", "result", "result save directory")
	return downloadCmd
}
