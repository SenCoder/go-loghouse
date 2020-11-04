package main

import (
	"flag"
	"github.com/sencoder/go-loghouse/config"
	"github.com/sencoder/go-loghouse/httpd"
	"github.com/sencoder/go-loghouse/pkg/log"
	"github.com/sirupsen/logrus"
)

var (
	configPath = flag.String("config", "", "configure file path")
)

func main() {
	flag.Parse()

	var cfg *config.Config
	var err error

	if *configPath != "" {
		cfg, err = config.LoadConfig(*configPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cfg = config.DefaultConfig
	}

	log.SetLevel(logrus.DebugLevel)
	log.Fatal(httpd.RunServer(cfg.Http))

}
