package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/* Transaction */
type Transaction struct {
	Txid     string `json:"txid"`
	Size     uint   `json:"size"`
	Version  uint   `json:"version"`
	Locktime uint   `json:"locktime"`
	Fee      uint   `json:"fee"`
	Inputs   []struct {
		Coinbase  bool     `json:"coinbase"`
		Txid      string   `json:"txid"`
		Output    uint     `json:"output"`
		Sigscript string   `json:"sigscript"`
		Sequence  uint64   `json:"sequence"`
		Pkscript  string   `json:"pkscript"`
		Value     uint     `json:"value"`
		Address   string   `json:"address"`
		Witness   []string `json:"witness"`
	} `json:"inputs"`
	Outputs []struct {
		Address  string `json:"address"`
		Pkscript string `json:"pkscript"`
		Value    uint   `json:"value"`
		Spent    bool   `json:"spent"`
		Spender  struct {
			Txid  string `json:"txid"`
			Input uint   `json:"input"`
		} `json:"spender,omitempty"`
		Input uint `json:"input,omitempty"`
	} `json:"outputs"`
	Block struct {
		Height   uint `json:"height"`
		Position uint `json:"position"`
	} `json:"block"`
	Deleted bool `json:"deleted"`
	Time    uint `json:"time"`
	Rbf     bool `json:"rbf"`
	Weight  uint `json:"weight"`
}

func GetTxInAddrs(tx Transaction) []string {
	var in_addrs []string
	for _, utxo := range tx.Inputs {
		if utxo.Address != "" {
			in_addrs = append(in_addrs, utxo.Address)
		}
	}
	return in_addrs
}

func GetTxOutAddrs(tx Transaction) []string {
	var out_addrs []string
	for _, utxo := range tx.Outputs {
		if utxo.Address != "" {
			out_addrs = append(out_addrs, utxo.Address)
		}
	}
	return out_addrs
}

func GetTxTime(tx Transaction) string {
	timeLayout := "2006-01-02 15:04:05"
	return time.Unix(int64(tx.Time), 0).Format(timeLayout)
}

// if given addr in tx inputs
func IsInTxInputs(addr string, tx Transaction) bool {
	in_addrs := GetTxInAddrs(tx)
	for _, cur_addr := range in_addrs {
		if cur_addr == addr {
			return true
		}
	}
	return false
}

// if given addr in tx outputs
func IsInTxOutputs(addr string, tx Transaction) bool {
	out_addrs := GetTxOutAddrs(tx)
	for _, cur_addr := range out_addrs {
		if cur_addr == addr {
			return true
		}
	}
	return false
}

func IsCoinbaseTx(tx Transaction) bool {
	return tx.Inputs[0].Coinbase
}

func ReadTransaction(path string) ([]Transaction, error) {
	var txs []Transaction
	obj, err := os.ReadFile(path)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("read file: %s error", path))
		return nil, err
	}
	err = json.Unmarshal(obj, &txs)
	if err != nil {
		err = errors.Wrap(err, "unmarshall error")
		return nil, err
	}
	return txs, nil
}

func ReadTransactionDir(blockDir string) ([]Transaction, error) {
	var all_txs []Transaction
	files, err := os.ReadDir(blockDir)
	if err != nil {
		log.Fatal(err)
	}
	n := len(files)
	bar := GetProgressBar(n)
	defer bar.Close()
	for _, file := range files {
		block_height_path := filepath.Join(blockDir, file.Name())
		bar.Describe(fmt.Sprintf("loading tx in %s:", file.Name()))
		txs, err := ReadTransaction(block_height_path)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("read %s error", block_height_path))
			return nil, err
		}
		all_txs = append(all_txs, txs...)
		bar.Add(1)
	}
	return all_txs, nil
}

func MashalTransactions(txDir string, savepath string) error {
	txs, err := ReadTransactionDir(txDir)
	if err != nil {
		err = errors.Wrap(err, "")
		return err
	}
	obj, err := json.Marshal(txs)
	if err != nil {
		err = errors.Wrap(err, "")
		return err
	}
	Save(savepath, obj, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	return nil
}

func ReadAddrs(path string) ([]string, error) {
	addrs := make([]string, 0)
	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("read %s failed", path))
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		addrs_str := strings.Split(line, " ")
		addrs = append(addrs, addrs_str...)
	}
	if err := scanner.Err(); err != nil {
		err = errors.Wrap(err, "scanner error")
		return nil, err
	}
	return addrs, nil
}
