package main

import (
	"flag"
	"log"

	"github.com/ChernichenkoStephan/nanostats/internal/tg"
	"go.uber.org/zap"
)

func main() {
	var confPath string
	flag.StringVar(&confPath,
		`config`, `config.yaml`,
		`path to config file`,
	)
	flag.Parse()

	cfg, err := readConfig(confPath)
	if err != nil {
		log.Fatal(cfg)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()
	lg := logger.Sugar()

	repo := tg.NewRepository()
	initRepository(cfg, repo)

	b, err := initBot(cfg, lg, repo)
	if err != nil {
		log.Fatal(err)
	}

	lg.Infoln("Setup done.")
	b.Start()
}
