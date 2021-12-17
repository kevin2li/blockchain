package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type BTCBlock struct {
	Hash      string   `json:"hash"`
	Height    uint     `json:"height"`
	Mainchain bool     `json:"mainchain"`
	Previous  string   `json:"previous"`
	Time      uint     `json:"time"`
	Version   uint     `json:"version"`
	Bits      uint     `json:"bits"`
	Nonce     uint64   `json:"nonce"`
	Size      uint     `json:"size"`
	Tx        []string `json:"tx"`
	Merkle    string   `json:"merkle"`
	Subsidy   uint     `json:"subsidy"`
	Fees      uint     `json:"fees"`
	Outputs   uint64   `json:"outputs"`
	// Work      uint64   `json:"work"`
	Weight uint `json:"weight"`
}

type ETHBlock struct {
	BlockHeaders []struct {
		Hash                     string        `json:"hash"`
		Number                   string        `json:"number"`
		ParentHash               string        `json:"parentHash"`
		Uncles                   []interface{} `json:"uncles"`
		Sha3Uncles               string        `json:"sha3Uncles"`
		TransactionsRoot         string        `json:"transactionsRoot"`
		StateRoot                string        `json:"stateRoot"`
		LogsBloom                string        `json:"logsBloom"`
		Difficulty               string        `json:"difficulty"`
		TotalDifficulty          string        `json:"totalDifficulty"`
		GasLimit                 string        `json:"gasLimit"`
		GasUsed                  string        `json:"gasUsed"`
		ExtraData                string        `json:"extraData"`
		Timestamp                string        `json:"timestamp"`
		Size                     string        `json:"size"`
		Miner                    string        `json:"miner"`
		Nonce                    string        `json:"nonce"`
		StaticReward             string        `json:"staticReward"`
		BlockReward              string        `json:"blockReward"`
		TotalUncleReward         string        `json:"totalUncleReward"`
		TotalFees                string        `json:"totalFees"`
		TransactionCount         int           `json:"transactionCount"`
		InternalTransactionCount int           `json:"internalTransactionCount"`
	} `json:"blockHeaders"`
	From string `json:"from"`
	Size int    `json:"size"`
}

func ReadHeights(path string) ([]int, error) {
	heights := make([]int, 0)
	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("read %s failed", path))
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		heights_str := strings.Split(line, " ")
		temp_heights, _ := Strings2Ints(heights_str)
		heights = append(heights, temp_heights...)
	}
	if err := scanner.Err(); err != nil {
		err = errors.Wrap(err, "scanner error")
		return nil, err
	}
	return heights, nil
}
