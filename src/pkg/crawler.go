package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Crawler struct {
	UAPool  []string // browser user agent pool
	Savedir string   // result save directory
}

type CrawlAPI interface {
	GetBlocksByHeights(heights string) ([]interface{}, error)
	GetTxsByHashs(txHashs string) ([]interface{}, error)
	DownloadAllBlocks(blocks []interface{})
}

type BTCCrawler struct {
	Crawler
	GetBlockUrl string // https://api.blockchain.info/haskoin-store/btc/block/heights?heights=%s&notx=false
	GetTxUrl    string // https://api.blockchain.info/haskoin-store/btc/transactions?txids=%s
	Size        int    // number of tx each request get
}

type ETHCrawler struct {
	Crawler
	GetBlockHashUrl string // https://api.blockchain.info/v2/eth/data/blocks?size=%d&from=%s
	GetTxHashUrl    string // https://api.blockchain.info/v2/eth/data/block/hash/%s
	Size        int    // number of tx each request get
}

func NewBTCCrawler(getBlockUrl string, getTxUrl string, savedir string, size int, ua_path string) (*BTCCrawler, error) {
	var obj = BTCCrawler{
		GetBlockUrl: getBlockUrl,
		GetTxUrl:    getTxUrl,
		Crawler: Crawler{
			Savedir: savedir,
			UAPool: nil,
		},
		Size: size,
	}
	if ua_path != "" {
		content, err := os.ReadFile(ua_path)
		if err != nil {
			err = errors.Wrap(err, "read ua_path failed")
			return nil, err
		}
		err = json.Unmarshal(content, &obj.UAPool)
		if err != nil {
			err = errors.Wrap(err, "unmarshall ua_pool failed")
			return nil, err
		}
	}
	return &obj, nil
}

func NewETHCrawler(getBlockUrl string, getTxUrl string, savedir string, size int, ua_path string) (*BTCCrawler, error) {
	var obj = BTCCrawler{
		GetBlockUrl: getBlockUrl,
		GetTxUrl:    getTxUrl,
		Crawler: Crawler{
			Savedir: savedir,
			UAPool: nil,
		},
		Size: size,
	}
	if ua_path != "" {
		content, err := os.ReadFile(ua_path)
		if err != nil {
			err = errors.Wrap(err, "read ua_path failed")
			return nil, err
		}
		err = json.Unmarshal(content, &obj.UAPool)
		if err != nil {
			err = errors.Wrap(err, "unmarshall ua_pool failed")
			return nil, err
		}
	}
	return &obj, nil
}

func (c *BTCCrawler) GetBlocksByHeights(heights string) ([]BTCBlock, error) {
	/* construct request */
	url := fmt.Sprintf(c.GetBlockUrl, heights)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", c.UAPool[rand.Intn(len(c.UAPool))])
	req.Header.Set("Content-Type", "application/json")

	/* issue request and wait response*/
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("request for %s failed!", url))
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = errors.Wrap(err, "read response failed")
		return nil, err
	}

	/* save response */
	var blocks []BTCBlock
	err = json.Unmarshal(body, &blocks)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("unmarshall error, response is:\n %s", string(body)[:200]))
		return nil, err
	}
	return blocks, nil
}

func (c *BTCCrawler) GetTxsByHashs(txHashs string) ([]Transaction, error) {
	/* construct request */
	url := fmt.Sprintf(c.GetTxUrl, txHashs)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("request for %s error!", url))
		return nil, err
	}
	req.Header.Set("User-Agent", c.UAPool[rand.Intn(len(c.UAPool))])
	req.Header.Set("Content-Type", "application/json")

	/* issue request and wait response*/
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("request for %s failed!", url))
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("read response failed, request url is: %s", url))
		return nil, err
	}

	/* save response */
	var txs []Transaction
	err = json.Unmarshal(body, &txs)
	if err != nil {
		err = errors.Wrap(err, "unmarshall failed\n")
		Save("error_page.html", body, os.O_CREATE|os.O_WRONLY)
		return nil, err
	}
	return txs, nil
}

func (c *BTCCrawler) GetBlocks(heights []int) ([]BTCBlock, error) {
	fmt.Printf("INFO: get txids in blocks with heights = %v...\n", heights)
	var all_blocks []BTCBlock
	var n = len(heights)
	bar := GetProgressBar(n)
	defer bar.Close()

	for _, h := range heights {
		bar.Describe(fmt.Sprintf("download txids in block %d :", h))
		blocks, err := c.GetBlocksByHeights(strconv.Itoa(h))
		if err != nil {
			return nil, err
		}
		all_blocks = append(all_blocks, blocks...)
		bar.Add(1)
	}
	return all_blocks, nil
}

func (c *BTCCrawler) GetBlocksInRange(low, high int) ([]BTCBlock, error) {
	fmt.Printf("INFO: get txids in blocks with heights in range [%d, %d)...\n", low, high)
	var all_blocks []BTCBlock
	var n = high - low
	bar := GetProgressBar(n)
	defer bar.Close()

	for i := low; i < high; i++ {
		bar.Describe(fmt.Sprintf("downloading txids in block %d :", i))
		blocks, err := c.GetBlocksByHeights(strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		all_blocks = append(all_blocks, blocks...)
		bar.Add(1)
	}
	return all_blocks, nil
}

func (c *BTCCrawler) DownloadOneBlock(block *BTCBlock, done chan int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v\n", err)
			done <- int(block.Height)
		} else {
			done <- 0
		}
	}()
	fmt.Printf("INFO: Downloading block at height %d...\n", block.Height)
	p, n, tx_hashs := 0, len(block.Tx), ""
	bar := GetProgressBar(n)
	defer bar.Close()
	bar.Describe(fmt.Sprintf("download block %d:", block.Height))
	var all_txs []Transaction
	for i, tx_hash := range block.Tx {
		tx_hashs = fmt.Sprintf("%s,%s", tx_hashs, tx_hash)
		p++
		// every `Size` hash issue a request
		if (p+1)%c.Size == 0 || i == n-1 {
			retry_times := 0
		request:
			txs, err := c.GetTxsByHashs(tx_hashs[1:])
			if err != nil {
				if retry_times < 5 {
					retry_times++
					time.Sleep(time.Duration(retry_times*10) * time.Second)
					fmt.Printf("Retry download block %d for %d time(s)...\n", block.Height, retry_times)
					goto request
				}
				panic(err)
			}
			all_txs = append(all_txs, txs...)
			bar.Add(p)
			p, tx_hashs = 0, ""
		}
	}
	obj, err := json.Marshal(all_txs)
	if err != nil {
		err = errors.Wrap(err, "Marshal Error")
		panic(err)
	}
	savepath := filepath.Join(c.Savedir, fmt.Sprintf("block_height=%d.json", block.Height))
	Save(savepath, obj, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	fmt.Printf("INFO: Block %d download success!\n", block.Height)
}

func (c *BTCCrawler) DownloadAllBlocks(blocks []BTCBlock) {
	n := len(blocks)
	failedBlocks := make([]int, 0)
	done := make(chan int, n)

	for i := 0; i < n; i++ {
		go c.DownloadOneBlock(&blocks[i], done)
		// decrease access frequency
		r := 60 + rand.Intn(5)
		time.Sleep(time.Duration(r) * time.Second)
	}
	for i := 0; i < n; i++ {
		if h := <-done; h != 0 {
			failedBlocks = append(failedBlocks, h)
		}
	}

	close(done)
	log.Printf("Total : %d, Success: %d, Failure: %d\n", n, n-len(failedBlocks), len(failedBlocks))
	if len(failedBlocks) > 0 {
		sort.Ints(failedBlocks)
		log.Printf("Failed blocks are: %v\n", failedBlocks)
		content := fmt.Sprintf("%v\n", failedBlocks)
		content = content[1:len(content)-2] + "\n"
		Save("failed_block_heights.txt", []byte(content), os.O_CREATE|os.O_WRONLY|os.O_APPEND)
		log.Printf("Save failed block heights at: %s", "failed_block_heights.txt")
	}
}

func (e *ETHCrawler) GetBlockHash(start_height string, size int) ([]ETHBlock, error) {
	/* construct request */
	url := fmt.Sprintf(e.GetBlockUrl, size, start_height)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", e.UAPool[rand.Intn(len(e.UAPool))])
	req.Header.Set("Content-Type", "application/json")

	/* issue request and wait response*/
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("request for %s failed!", url))
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = errors.Wrap(err, "read response failed")
		return nil, err
	}

	/* save response */
	var blocks []ETHBlock
	err = json.Unmarshal(body, &blocks)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("unmarshall error, response is:\n %s", string(body)[:200]))
		return nil, err
	}
	return blocks, nil
}