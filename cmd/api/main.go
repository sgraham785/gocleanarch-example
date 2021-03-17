package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
	bookAdapter "github.com/sgraham785/gocleanarch-example/internal/book/adapter"
	bookInfra "github.com/sgraham785/gocleanarch-example/internal/book/infrastructure"
	bookUseCase "github.com/sgraham785/gocleanarch-example/internal/book/usecase"
	borrowAdapter "github.com/sgraham785/gocleanarch-example/internal/borrow/adapter"
	borrowUseCase "github.com/sgraham785/gocleanarch-example/internal/borrow/usecase"
	userAdapter "github.com/sgraham785/gocleanarch-example/internal/user/adapter"
	userInfra "github.com/sgraham785/gocleanarch-example/internal/user/infrastructure"
	userUseCase "github.com/sgraham785/gocleanarch-example/internal/user/usecase"
	"github.com/sgraham785/gocleanarch-example/pkg/config"
	"github.com/sgraham785/gocleanarch-example/pkg/logger"
	"github.com/sgraham785/gocleanarch-example/pkg/repository"
	"github.com/sgraham785/gocleanarch-example/pkg/router"
	"github.com/sgraham785/gocleanarch-example/pkg/server"
)

func main() {
	cfg := config.Load()

	logger := logger.New()
	defer logger.Zap.Sync()

	db := repository.NewPostgresConn(cfg.PostgresConf)
	defer db.Pg.Close()

	r := router.NewChiRouter()

	server := &server.Server{
		Cfg:    cfg,
		Log:    logger,
		DB:     db,
		Router: r,
	}

	bookRepo := bookInfra.NewPgRepo(server)
	bookUseCase := bookUseCase.New(server, bookRepo)

	userRepo := userInfra.NewPgRepo(server)
	userUseCase := userUseCase.New(server, userRepo)
	borrowUseCase := borrowUseCase.New(server, userUseCase, bookUseCase)

	bookAdapter.HTTPRoutes(server, bookUseCase)
	userAdapter.HTTPRoutes(server, userUseCase)
	borrowAdapter.HTTPRoutes(server, bookUseCase, userUseCase, borrowUseCase)

	r.Chi.Handle("/", r.Chi)
	r.Chi.Handle("/metrics", promhttp.Handler())
	r.Chi.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(cfg.APIPort),
		Handler:      r.Chi,
	}
	server.Log.Zap.Info("Started on -> ", zap.Int("port", cfg.APIPort))
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
