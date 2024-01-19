package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ruhollahh/paperback/internal/app/domain"
	"time"
)

type ProductService struct {
	DB *sql.DB
}

func (s ProductService) GetAll(title string, filters Filters) ([]domain.Product, Metadata, error) {
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, created_at, title, description, price, version
        FROM products
        WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')   
        ORDER BY %s %s, id ASC
        LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, filters.limit(), filters.offset()}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	products := []domain.Product{}
	totalRecords := 0

	for rows.Next() {
		var product domain.Product

		err := rows.Scan(
			&totalRecords,
			&product.ID,
			&product.CreatedAt,
			&product.Title,
			&product.Description,
			&product.Price,
			&product.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := NewMetadata(totalRecords, filters.Page, filters.PageSize)

	return products, metadata, nil
}

type CreateProductReq struct {
	Title       string
	Description string
	Price       int32
}

type CreateProductRes struct {
	ID        int64
	CreatedAt time.Time
	Version   int32
}

func (s ProductService) CreateProduct(req CreateProductReq) (*CreateProductRes, error) {
	//todo: validate the input
	query := `
        INSERT INTO products (title, description, price) 
        VALUES ($1, $2, $3)
        RETURNING id, created_at, version`

	args := []any{req.Title, req.Description, req.Price}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res CreateProductRes
	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&res.ID, &res.CreatedAt, &res.Version)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (s ProductService) Get(id int64) (*domain.Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, title, description, price, version
        FROM products
        WHERE id = $1`

	var product domain.Product

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.Title,
		&product.Description,
		&product.Price,
		&product.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

type UpdateProductReq struct {
	ID          int64
	Title       string
	Description string
	Price       int32
	Version     int32
}

type UpdateProductRes struct {
	Version int32
}

func (s ProductService) Update(req UpdateProductReq) (*UpdateProductRes, error) {
	// todo: validate the input
	query := `
        UPDATE products 
        SET title = $1, description = $2, price = $3, version = version + 1
        WHERE id = $4 AND version = $5
        RETURNING version`

	args := []any{
		req.Title,
		req.Description,
		req.Price,
		req.ID,
		req.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res UpdateProductRes
	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&res.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrEditConflict
		default:
			return nil, err
		}
	}

	return &res, nil
}

func (s ProductService) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM products
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
