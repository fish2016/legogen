package utils

import (
	"bytes"
	"github.com/fish2016/legogen/logger"
	"go/format"
	"log"
	"os"
	"path/filepath"
)

func FormatBuffer(buf bytes.Buffer, filename string) []byte {
	src, err := format.Source(buf.Bytes())
	if err != nil {
		logger.Logger.Printf("Warning: internal Error: invalid go generated in file %s: %s", filename, err)
		logger.Logger.Printf("Warning: compile the package to analyze the error: %s", err)
		return buf.Bytes()
	}
	return src
}

func OpenFile(dirname, filename string) *os.File {
	_, err := os.Stat(dirname)
	err = os.MkdirAll(dirname, 0744)
	if err != nil && !os.IsExist(err) {
		logger.Logger.Fatalf("Unable to Make Directory: %s: %s", dirname, err)
	}
	fname := filepath.Join(dirname, filename)

	file, err := os.Create(fname)
	if err != nil && os.IsExist(err) {
		file, err = os.Open(fname)
	}

	if err != nil {
		log.Fatalf("Unable to open or create file %s: %s", fname, err)
	}
	return file
}

// ChecExist check if the file exists
func CheckExist(src string) bool {
	_, err := os.Stat(src)

	return !os.IsNotExist(err)
}
