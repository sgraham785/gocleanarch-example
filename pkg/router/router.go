package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	mw "github.com/sgraham785/gocleanarch-example/pkg/middleware"
)

type HTTPRouter struct {
	Chi *chi.Mux
}

func NewChiRouter() *HTTPRouter {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(mw.NewMetrics("goarch"))

	return &HTTPRouter{Chi: r}
}
