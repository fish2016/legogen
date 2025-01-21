package utils

import (
	"log"
	"os"
	"path/filepath"
)

func IsDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func PrefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}

	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}
