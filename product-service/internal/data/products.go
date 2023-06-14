package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var ErrRecordNotFound = errors.New("record not found")

type ProductModel struct {
	DB *sql.DB
}

func (p ProductModel) Insert(product *proto.Product) (*proto.Product, error) {
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

	var creationDate time.Time
	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&product.Id, &creationDate, &product.Version)
	if err != nil {
		return nil, err
	}
	product.CreationDate = timestamppb.New(creationDate)

	return product, nil
}

func (p ProductModel) Get(id int64) (*proto.Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, name, price, description, category, is_available, creation_date, version
			  FROM products
		      WHERE id = $1`

	var product proto.Product
	var creationDate time.Time
	err := p.DB.QueryRow(query, id).Scan(
		&product.Id,
		&product.Name,
		&product.Price,
		&product.Description,
		&product.Category,
		&product.IsAvailable,
		&creationDate,
		&product.Version,
	)

	product.CreationDate = timestamppb.New(creationDate)

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

func (p ProductModel) GetAll(name string, category string, filters *proto.Filters) ([]*proto.Product, *proto.Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, name, price, description, category, is_available, creation_date, version
		FROM products
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', category) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, sortColumn(filters), sortDirection(filters))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []any{name, category, limit(filters), offset(filters)}

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, &proto.Metadata{}, err
	}

	defer rows.Close()

	var totalRecords int32 = 0

	var products []*proto.Product

	for rows.Next() {
		var product proto.Product
		var creationDate time.Time
		err := rows.Scan(
			&totalRecords,
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Category,
			&product.IsAvailable,
			&creationDate,
			&product.Version,
		)
		product.CreationDate = timestamppb.New(creationDate)
		if err != nil {
			return nil, &proto.Metadata{}, err
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, &proto.Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters)
	return products, metadata, nil
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
