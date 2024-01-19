package api

import (
	"log/slog"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/ruhollahh/paperback/api/config"
	"github.com/ruhollahh/paperback/internal/app/service"
	"github.com/ruhollahh/paperback/internal/mailer"
)

type API struct {
	Logger         *slog.Logger
	Config         config.Config
	Services       service.Services
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
	Mailer         mailer.Mailer
	Wg             sync.WaitGroup
}
