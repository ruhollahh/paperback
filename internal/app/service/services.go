package service

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/ruhollahh/paperback/api/config"
)

type Services struct {
	Tokens      TokenService
	Users       UserService
	Permissions PermissionsService
	Products    ProductService
}

func NewServices(db *sql.DB) Services {
	return Services{
		Tokens:      TokenService{DB: db},
		Users:       UserService{DB: db},
		Permissions: PermissionsService{DB: db},
		Products:    ProductService{DB: db},
	}
}

func OpenDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)

	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)

	db.SetConnMaxIdleTime(cfg.Db.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
