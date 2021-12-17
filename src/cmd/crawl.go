package cmd

import (
	"blockchain/src/pkg"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type DataSource struct {
	BTC struct {
		GetBlockURL string `json:"GetBlockUrl"`
		GetTxURL    string `json:"GetTxUrl"`
	} `json:"BTC"`
	ETH struct {
		GetBlockURL string `json:"GetBlockUrl"`
		GetTxURL    string `json:"GetTxUrl"`
	} `json:"ETH"`
}

func CrawlSetup() *cobra.Command {

	var (
		isInterval bool
		filepath   string
		platform   string
		savedir    string
		size       int
		ua_path    string
	)
	var data_source_path = "/home/likai/code/go_program/go_learn2/src/config/data_source.json"
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
			content, err := os.ReadFile(data_source_path)
			if err != nil {
				err = errors.Wrap(err, fmt.Sprintf("read %s failed", data_source_path))
				log.Fatal(err)
			}
			var d DataSource
			json.Unmarshal(content, &d)
			var getBlockUrl, getTxUrl string
			switch platform {
			case "BTC":
				getBlockUrl = d.BTC.GetBlockURL
				getTxUrl = d.BTC.GetTxURL
			case "ETH":
				getBlockUrl = d.ETH.GetBlockURL
				getTxUrl = d.ETH.GetTxURL
			}

			crawler, _ := pkg.NewCrawler(
				platform,
				getBlockUrl,
				getTxUrl,
				savedir,
				size,
				ua_path,
			)
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

	downloadCmd.Flags().StringVarP(&platform, "platform", "p", "BTC", "BTC and ETH are  now supported platform")
	downloadCmd.Flags().BoolVarP(&isInterval, "interval", "r", false, "block heights range")
	downloadCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "file store heights to download")
	downloadCmd.Flags().StringVarP(&savedir, "savedir", "s", "result", "result save directory")
	downloadCmd.Flags().StringVar(&ua_path, "ua_path", "/home/likai/code/go_program/go_learn2/src/config/ua.json", "User-Agent pool")
	downloadCmd.Flags().IntVar(&size, "size", 10, "request number of txs each time")
	downloadCmd.MarkFlagRequired("platform")

	return downloadCmd
}
