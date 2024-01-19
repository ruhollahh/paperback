package service

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
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

type DBConfig struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

func OpenDB(cfg DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)

	db.SetMaxIdleConns(cfg.MaxIdleConns)

	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
