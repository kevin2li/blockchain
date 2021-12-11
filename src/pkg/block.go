package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Block struct {
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
