package generator

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"path/filepath"

	. "github.com/fish2016/legogen/logger"
)

func sliceContains(slice []string, entry string) bool {
	for _, s := range slice {
		if entry == s {
			return true
		}
	}
	return false
}

func DetermineLocalName(suggestedName string, currentNames []string) string {
	if !sliceContains(currentNames, suggestedName) {
		return suggestedName
	}

	f := func(prefix string) string {
		p := []byte(prefix)
		for i := 97; i < 97+26; i++ {
			b := append([]byte{}, p...)
			b = append(b, byte(i))
			if !sliceContains(currentNames, string(b)) {
				return string(b)
			}
		}
		return ""
	}

	for sliceContains(currentNames, suggestedName) {
		suggestedName = suggestedName + "a"

		suggestedName = f(suggestedName)
	}

	return suggestedName
}

func FormatBuffer(buf bytes.Buffer, filename string) []byte {
	src, err := format.Source(buf.Bytes())
	if err != nil {
		Logger.Printf("Warning: internal Error: invalid go generated in file %s: %s", filename, err)
		Logger.Printf("Warning: compile the package to analyze the error: %s", err)
		return buf.Bytes()
	}
	return src
}

func OpenFile(dirname, filename string) *os.File {
	_, err := os.Stat(dirname)
	err = os.MkdirAll(dirname, 0744)
	if err != nil && !os.IsExist(err) {
		Logger.Fatalf("Unable to Make Directory: %s: %s", dirname, err)
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

func IsDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func prefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}

	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}
