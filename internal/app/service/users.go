package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ruhollahh/paperback/pkg/errsx"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/ruhollahh/paperback/internal/app/domain"
)

type UserService struct {
	DB *sql.DB
}

var AnonymousUser = &domain.User{}

func IsAnonymous(u *domain.User) bool {
	return u == AnonymousUser
}

type password struct {
	plaintext string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type SignupReq struct {
	Name     string
	Email    string
	Password string
}

type SignupRes struct {
	ID        int64
	CreatedAt time.Time
	Version   int32
}

func (s UserService) Signup(req SignupReq) (*SignupRes, error) {
	var input struct {
		name     string
		email    string
		password password
	}
	var err error
	var errs errsx.Map

	input.name, err = domain.NewName(req.Name)
	if err != nil {
		errs.Set("name", err)
	}
	input.email, err = domain.NewEmail(req.Email)
	if err != nil {
		errs.Set("email", err)
	}
	plaintextPassword, err := domain.NewEmail(req.Password)
	if err != nil {
		errs.Set("password", err)
	}
	if errs != nil {
		return nil, fmt.Errorf("%w: %w", ErrBadRequest, errs)
	}

	err = input.password.Set(plaintextPassword)
	if err != nil {
		return nil, err
	}

	query := `
        INSERT INTO users (name, email, password_hash, activated) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

	args := []any{input.name, input.email, input.password.hash, false}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res SignupRes
	err = s.DB.QueryRowContext(ctx, query, args...).Scan(&res.ID, &res.CreatedAt, &res.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	return &res, nil
}

func (s UserService) GetByID(id int64) (*domain.User, error) {
	query := `
        SELECT id, created_at, name, email, password_hash, activated, version
        FROM users
        WHERE id = $1`

	var user domain.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s UserService) GetByEmail(email string) (*domain.User, error) {
	_, err := domain.NewEmail(email)
	if err != nil {
		var errs errsx.Map
		errs.Set("email", err)
		return nil, errs
	}

	query := `
        SELECT id, created_at, name, email, password_hash, activated, version
        FROM users
        WHERE email = $1`

	var user domain.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = s.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

type ActivateUserReq struct {
	ID      string
	Version int32
}

func (s UserService) ActivateUser(req ActivateUserReq) error {
	query := `
        UPDATE users 
        SET activated = $1, version = version + 1
        WHERE id = $2 AND version = $3
        RETURNING version`

	args := []any{
		req.ID,
		req.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&req.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

type GetForTokenReq struct {
	TokenScope     domain.TokenScope
	TokenPlaintext string
}

func (s UserService) GetForToken(req GetForTokenReq) (*domain.User, error) { /**/
	tokenHash := sha256.Sum256([]byte(req.TokenPlaintext))

	query := `
        SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = $1
        AND tokens.scope = $2 
        AND tokens.expiry > $3`

	args := []any{tokenHash[:], req.TokenScope, time.Now()}

	var user domain.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
