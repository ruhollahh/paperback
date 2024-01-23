package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/ruhollahh/paperback/api/config"
	"github.com/ruhollahh/paperback/api/handler"
	"github.com/ruhollahh/paperback/api/middleware"
)

func (a *API) routes() http.Handler {
	router := httprouter.New()

	handler := &handler.Handler{
		Router:         router,
		Config:         a.Config,
		Version:        config.Version,
		Logger:         a.Logger,
		Services:       a.Services,
		FormDecoder:    a.FormDecoder,
		SessionManager: a.SessionManager,
	}

	middleware := &middleware.Middleware{
		Logger:         a.Logger,
		Services:       a.Services,
		FormDecoder:    a.FormDecoder,
		SessionManager: a.SessionManager,
	}

	// fileServer := http.FileServer(http.FS(web.Files))
	// router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	fileServer := http.FileServer(http.Dir("./web/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(a.SessionManager.LoadAndSave, middleware.CSRF, middleware.Authenticate)
	// protected := dynamic.Append(app.requireAuth)

	router.HandlerFunc(http.MethodGet, "/healthcheck", handler.HealthCheck)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(handler.Home))

	router.NotFound = http.HandlerFunc(handler.NotFound)

	standard := alice.New(middleware.RecoverPanic, middleware.LogRequest, middleware.SecureHeaders)
	return standard.Then(router)
}
