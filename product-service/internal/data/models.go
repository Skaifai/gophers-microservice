package data

import (
	"database/sql"
)

type Models struct {
	Products ProductModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Products: ProductModel{DB: db},
	}
}
