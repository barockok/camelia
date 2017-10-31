package main

import (
	"bytes"
	"flag"

	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "location of config file")
	flag.Parse()
	loadConfig(configPath)
}

func loadConfig(configPath string) {
	if configPath == "" {
		logrus.Fatalf("config is required argument")
	}

	cfgCont, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Fatalf("Cannot read config file")
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(cfgCont))
	if err != nil {
		logrus.Fatalf("Cannot parse configuration, err : %v", err)
	}
}
