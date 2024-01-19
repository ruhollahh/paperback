package main

import (
	"flag"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/ruhollahh/paperback/api"
	"github.com/ruhollahh/paperback/api/config"
	"github.com/ruhollahh/paperback/internal/app/service"
	"github.com/ruhollahh/paperback/internal/mailer"
)

func main() {
	var cfg config.Config

	flag.IntVar(&cfg.Port, "port", 4000, "Web server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.Db.Dsn, "db-dsn", os.Getenv("PAPERBACK_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.Smtp.Host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", "b316465853018d", "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", "bdf309a7f75e60", "SMTP password")
	flag.StringVar(&cfg.Smtp.Sender, "smtp-sender", "Paperback <no-reply@paperback.com>", "SMTP sender")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.Cors.TrustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	db, err := service.OpenDB(cfg.Db)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	a := &api.API{
		Config:         cfg,
		Logger:         logger,
		Services:       service.NewServices(db),
		Mailer:         mailer.NewMailer(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
	}

	err = a.Serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
