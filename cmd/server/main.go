package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	bookUseCase "github.com/sgraham785/gocleanarch-example/internal/book/usecase"

	_ "github.com/lib/pq"
	"github.com/sgraham785/gocleanarch-example/config"

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
	c := config.Load()
	metricService, err := metric.NewPrometheusService(c)
	if err != nil {
		log.Fatal(err.Error())
	}
	appMetric := metric.NewCLI("search")
	appMetric.Started()
	query, err := handleParams()
	if err != nil {
		log.Fatal(err.Error())
	}

	db := repository.NewPostgresConn(c.PostgresConf)
	defer db.Close()
	repo := bookInfra.NewPgRepo(db)
	service := bookUseCase.New(repo)
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
