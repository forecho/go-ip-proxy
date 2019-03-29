package main

import (
	"fmt"
	"go-ip-proxy/logger"
	"go-ip-proxy/parser"
	"go-ip-proxy/storage"
	"go-ip-proxy/verifier"
	"go.uber.org/zap"
	"sync"
	"time"
)

func main() {
	log := logger.Config("./logs/all.log", "info", true)
	log.Info("test log", zap.Int("line", 47))

	configs := parser.NewConfig("./config/proxyWebsiteConfig.json")

	// Load database.
	database, err := storage.NewStorage()
	defer database.Close()
	if err != nil {
		logger.Error("db error", zap.Error(err))
		panic(err)
	}

	// Verify storage every 5min.
	verifyTicker := time.NewTicker(time.Minute * 5)
	go func() {
		for range verifyTicker.C {
			verifier.VerifyAndDelete(database)
			logger.Debug("verify database.")
		}
	}()

	run(database, configs)

}

func run(storage storage.Storage, configs *parser.Configs) {

	for {
		var wg sync.WaitGroup

		for _, configuration := range configs.Configs {
			items := parser.NewSelector(configuration)
			verifier.VerifyAndSave(items, storage)
		}
		for _, item := range storage.GetAll() {
			fmt.Printf("%s\n", string(item))
		}

		wg.Wait()
		logger.Debug("finish once, sleep 10 Second.")
		time.Sleep(time.Second * 10)
	}
}
