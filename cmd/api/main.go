package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sgraham785/gocleanarch-example/config"
	bookAdapter "github.com/sgraham785/gocleanarch-example/internal/book/adapter"
	bookInfra "github.com/sgraham785/gocleanarch-example/internal/book/infrastructure"
	bookUseCase "github.com/sgraham785/gocleanarch-example/internal/book/usecase"
	borrowAdapter "github.com/sgraham785/gocleanarch-example/internal/borrow/adapter"
	borrowUseCase "github.com/sgraham785/gocleanarch-example/internal/borrow/usecase"
	userAdapter "github.com/sgraham785/gocleanarch-example/internal/user/adapter"
	userInfra "github.com/sgraham785/gocleanarch-example/internal/user/infrastructure"
	userUseCase "github.com/sgraham785/gocleanarch-example/internal/user/usecase"
	"github.com/sgraham785/gocleanarch-example/pkg/metric"
	"github.com/sgraham785/gocleanarch-example/pkg/middleware"
	"github.com/sgraham785/gocleanarch-example/pkg/repository"
)

func main() {
	c := config.Load()
	db := repository.NewPostgresConn(c.PostgresConf)
	defer db.Close()

	bookRepo := bookInfra.NewPgRepo(db)
	bookUseCase := bookUseCase.New(bookRepo)

	userRepo := userInfra.NewPgRepo(db)
	userUseCase := userUseCase.New(userRepo)

	loanUseCase := borrowUseCase.New(userUseCase, bookUseCase)

	metricService, err := metric.NewPrometheusService(c)
	if err != nil {
		log.Fatal(err.Error())
	}
	r := mux.NewRouter()
	//handlers
	n := negroni.New(
		negroni.HandlerFunc(middleware.Cors),
		negroni.HandlerFunc(middleware.Metrics(metricService)),
		negroni.NewLogger(),
	)
	bookAdapter.BookRouter(r, *n, bookUseCase)
	userAdapter.UserRouter(r, *n, userUseCase)
	borrowAdapter.BorrowRouter(r, *n, bookUseCase, userUseCase, loanUseCase)

	http.Handle("/", r)
	http.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(c.APIPort),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
