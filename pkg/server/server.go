package server

import (
	"github.com/sgraham785/gocleanarch-example/pkg/config"
	"github.com/sgraham785/gocleanarch-example/pkg/logger"
	"github.com/sgraham785/gocleanarch-example/pkg/repository"
	"github.com/sgraham785/gocleanarch-example/pkg/router"
)

// Server struct is an object that represents the service, and holds all of its dependencies.
type Server struct {
	Cfg    *config.Specification
	Log    *logger.Logger
	DB     *repository.Repository
	Router *router.HTTPRouter
}
