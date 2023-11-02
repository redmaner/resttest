package main

import (
	"log"
	"os"

	"github.com/redmaner/resttest"
	"go.uber.org/zap"
)

func main() {

	logger, err := resttest.NewLogger()
	if err != nil {
		log.Fatalf("failed to create new logger: %s", err)
	}

	args := os.Args
	if len(args) < 2 {
		logger.Fatal("Expected argument")
	}

	configPath := args[1]

	config, err := resttest.LoadConfig(configPath)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	transport, _ := resttest.NewTransport(config.Headers)

	resttest.ExecuteTests(logger, transport, config.BaseUrl, config.Tests)
}
