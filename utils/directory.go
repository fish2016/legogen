package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func IsDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			wd, _ := os.Getwd()
			log.Printf("pwd: %s", wd)
			log.Fatal(errors.Join(fmt.Errorf("NOT_FOUND"), err))
		} else {
			log.Fatal(err)
		}

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
