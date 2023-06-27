package src

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

var logs = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func ReadYamlFromFile(configFile string, obj interface{}) error {
	defer func() error {
		if err := recover(); err != nil {
			logs.Fatal(fmt.Sprintf("read config file error %s", configFile))
			return errors.New(fmt.Sprintf("%v", err))
		}
		return nil
	}()

	// load yaml file
	if !Exists(configFile) {
		return nil
	}
	files, err := os.OpenFile(configFile, os.O_RDONLY, 0666)
	if err != nil {
		logs.Fatal(fmt.Sprintf("Read Laipvt Config File Error %v", err))
		return err
	}
	data, err := io.ReadAll(files)
	if err != nil {
		logs.Fatal(fmt.Sprintf("Read Laipvt Config File Error %v", err))
		return err
	}
	err = yaml.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	// print parsed configuration
	return nil
}
