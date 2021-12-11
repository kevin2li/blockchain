package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"

	pb "github.com/schollz/progressbar/v3"
)

func GetProgressBar(max int) *pb.ProgressBar {
	bar := pb.NewOptions(max,
		// pb.OptionSetWriter(ansi.NewAnsiStdout()),
		pb.OptionEnableColorCodes(true),
		pb.OptionShowBytes(true),
		pb.OptionSetWidth(15),
		pb.OptionShowCount(),
		pb.OptionThrottle(65*time.Millisecond),
		pb.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		pb.OptionSetTheme(pb.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	return bar
}

func Save(path string, content []byte, flag int) error {
	// check path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir := filepath.Dir(path)
		os.MkdirAll(dir, 0766)
	}
	// open or create file for writing
	file, err := os.OpenFile(path, flag, 0666)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Open file `%s` error", path))
		return err
	}
	defer file.Close()
	// write content
	writer := bufio.NewWriter(file)
	_, err = writer.Write(content)
	writer.Flush()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("save file `%s` error", path))
		return err
	}
	return nil
}

func Strings2Ints(strs []string) ([]int, error) {
	var result []int
	for _, s := range strs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, nil
}
