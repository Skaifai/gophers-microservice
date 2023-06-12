package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type ProductModel struct {
	DB *sql.DB
}

func (p ProductModel) Insert(product *proto.Product) (int64, error) {
	query := `INSERT INTO products (name, price, description, category, quantity, is_available)
			  VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id, creation_date, version`

	args := []any{
		product.Name,
		product.Price,
		product.Description,
		product.Category,
		product.Quantity,
		product.IsAvailable,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&product.Id, &product.CreationDate, &product.Version)

	if err != nil {
		return 0, err
	}

	return product.Id, nil
}

func (p ProductModel) Get(id int64) (*proto.Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, name, price, description, category, is_available, creation_date, version
			  FROM products
		      WHERE id = $1`

	var product proto.Product
	err := p.DB.QueryRow(query, id).Scan(
		&product.Id,
		&product.Name,
		&product.Price,
		&product.Description,
		&product.Category,
		&product.IsAvailable,
		&product.CreationDate,
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

func (p ProductModel) Update(product *proto.Product) error {
	query := `UPDATE products
	          SET name = $1, price = $2, description = $3, category = $4, quantity = $5, version = version + 1
	          WHERE id = $6
	          RETURNING version`

	args := []any{
		product.Name,
		product.Price,
		product.Description,
		product.Category,
		product.Quantity,
		product.Id,
	}

	return p.DB.QueryRow(query, args...).Scan(&product.Version)
}

func (p ProductModel) Delete(id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := p.DB.Exec(query, id)
	if err != nil {
		return nil
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
