package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	bookUseCase "github.com/sgraham785/gocleanarch-example/internal/book/usecase"

	_ "github.com/lib/pq"
	"github.com/sgraham785/gocleanarch-example/pkg/config"
	"github.com/sgraham785/gocleanarch-example/pkg/logger"
	"github.com/sgraham785/gocleanarch-example/pkg/server"

	bookInfra "github.com/sgraham785/gocleanarch-example/internal/book/infrastructure"
	"github.com/sgraham785/gocleanarch-example/pkg/metric"
	"github.com/sgraham785/gocleanarch-example/pkg/repository"
)

func handleParams() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("Invalid query")
	}
	return os.Args[1], nil
}

func main() {
	cfg := config.Load()

	logger := logger.New()
	defer logger.Zap.Sync()

	metricService, err := metric.NewPrometheusService(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	appMetric := metric.NewCLI("search")
	appMetric.Started()
	query, err := handleParams()
	if err != nil {
		log.Fatal(err.Error())
	}

	db := repository.NewPostgresConn(cfg.PostgresConf)
	defer db.Pg.Close()

	server := &server.Server{
		Cfg: cfg,
		Log: logger,
		DB:  db,
	}

	bookRepo := bookInfra.NewPgRepo(server)
	service := bookUseCase.New(server, bookRepo)
	all, err := service.SearchBooks(query)
	if err != nil {
		log.Fatal(err)
	}
	for _, j := range all {
		fmt.Printf("%s %s \n", j.Title, j.Author)
	}
	appMetric.Finished()
	err = metricService.SaveCLI(appMetric)
	if err != nil {
		log.Fatal(err)
	}
}
