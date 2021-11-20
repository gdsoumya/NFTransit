package main

import (
	"fmt"
	"github.com/gdsoumya/nftransit/relay/pkg/binding"
	"github.com/gdsoumya/nftransit/relay/pkg/database"
	"github.com/gdsoumya/nftransit/relay/pkg/env"
	"github.com/gdsoumya/nftransit/relay/pkg/handlers"
	"github.com/gdsoumya/nftransit/relay/pkg/queue"
	"github.com/gdsoumya/nftransit/relay/pkg/server"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"log"
	"math/big"
	"runtime"
)

func main() {
	var envConfig env.EnvData
	err := envconfig.Process("", &envConfig)
	if err != nil {
		log.Fatalf("failed to parse envs, err: %v", err.Error())
	}

	// logging level, dev mode enables debug logs
	var logger *zap.Logger

	// set log level
	if envConfig.Debug {
		logger, err = zap.NewDevelopment()
		logger.Debug("debug mode enabled")
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal("failed to create logger")
	}

	logger.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	logger.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))

	// create database connection
	dbClient, err := initDB(logger, &envConfig, false)
	if err != nil {
		logger.Fatal("failed to get db client", zap.Error(err))
	}

	// create queue clients
	// TODO(warning): amqp client might not be thread safe
	pubClient, err := queue.InitQueue(logger, &envConfig, true)
	if err != nil {
		logger.Fatal("failed to get publisher queue client", zap.Error(err))
	}
	defer pubClient.Close()

	//create chain clients
	evmClient, err := binding.NewEVMClient(logger, envConfig.NodeURL, envConfig.QueryPK, big.NewInt(envConfig.ChainID))
	if err != nil {
		logger.Fatal("failed to get chain client", zap.Error(err))
	}

	// init server handler
	handler := handlers.Handler{
		DB:        dbClient,
		QClient:   pubClient,
		Logger:    logger,
		EvmClient: evmClient,
	}

	block := false

	if envConfig.WorkerOn {
		handler.StartWorkers(envConfig)
		block = true
	}

	if envConfig.ServerOn {
		relaySrv := server.NewDefaultRelayServer(logger, handler, envConfig.Debug, envConfig.ServerPort)
		relaySrv.Init(true)
	} else if block {
		select {}
	}
}

func initDB(logger *zap.Logger, envConfig *env.EnvData, skipMigration bool) (*database.Database, error) {
	return database.NewDatabaseClient(logger, &database.DBConfig{
		User:     envConfig.DBUser,
		Password: envConfig.DBPassword,
		DBName:   envConfig.DBName,
		Host:     envConfig.DBHost,
		Port:     envConfig.DBPort,
	}, skipMigration)
}
