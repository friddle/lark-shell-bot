package src

import (
	"log"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

var logs = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
