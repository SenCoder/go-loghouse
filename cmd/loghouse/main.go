package main

import (
	"flag"
	"github.com/sencoder/go-loghouse/config"
	"github.com/sencoder/go-loghouse/httpd"
	"github.com/sencoder/go-loghouse/pkg/clickhouse"
	"github.com/sencoder/go-loghouse/pkg/log"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	configPath = flag.String("config", "", "configure file path.")
	deployDB   = flag.Bool("db-deploy", false, "run clickhouse database init only.")
	cluster    = flag.String("cluster", "", "clickhouse cluster name. cluster name is used for create tables in db-deploy.")
	ttl        = flag.Uint("ttl", 30, "clickhouse logs retention time in days. ttl is used for create tables in db-deploy.")
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

	if *deployDB {
		// init database
		ckInit := clickhouse.NewInitializer(cfg.Clickhouse, *cluster, *ttl)
		ckInit.CreateMergeTreeTable()
		ckInit.CreateBufferTable()
		ckInit.CreateDistributedTable()
		log.Info("clickhouse init success. ")
		os.Exit(0)
	}

	log.SetLevel(logrus.DebugLevel)
	log.Fatal(httpd.RunServer(cfg.Http))
}
