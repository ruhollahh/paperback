package handler

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/ruhollahh/paperback/api/config"
	"github.com/ruhollahh/paperback/internal/app/service"
)

type Handler struct {
	Router         *httprouter.Router
	Config         config.Config
	Version        string
	Logger         *slog.Logger
	Services       service.Services
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}
