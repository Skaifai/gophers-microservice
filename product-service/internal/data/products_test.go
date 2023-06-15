package data

import (
	"errors"
	"github.com/Skaifai/gophers-microservice/product-service/cmd/utils"
	"github.com/Skaifai/gophers-microservice/product-service/config"
	"github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"testing"
)

const id = 5

var products = func() ProductModel {
	cfg := loadTestConfiguration()
	db, _ := utils.OpenDB(cfg)
	return ProductModel{DB: db}
}()

func TestAddProduct(t *testing.T) {
	product := &proto.Product{
		Name:        "Apple",
		Price:       850,
		Description: "Apple from Almaty city",
		Category:    "Fruit",
		Quantity:    5,
	}
	SetStatus(product)
	if !product.IsAvailable {
		t.Error("product should be available")
	}

	_, err := products.Insert(product)
	if err != nil {
		t.Fatalf("error acquired while inserting product. %s", err.Error())
	}
}

func TestGetProduct(t *testing.T) {
	product, err := products.Get(id)
	if err != nil {
		t.Fatalf("error acquired while accessing table product. %s", err.Error())
	}
	t.Log(product)
}

func TestGetAllProduct(t *testing.T) {
	var input struct {
		Name     string
		Category string
		Filters  proto.Filters
	}
	input.Name = ""
	input.Category = ""
	input.Filters.Page = 1
	input.Filters.PageSize = 20
	input.Filters.Sort = "id"
	input.Filters.SortSafeList = []string{"id", "name", "category", "price", "is_available", "creation_date",
		"-id", "-name", "-category", "-price", "-is_available", "-creation_date"}

	products, metadata, err := products.GetAll(input.Name, input.Category, &input.Filters)
	if err != nil {
		t.Fatalf("error acquired while accessing table product. %s", err.Error())
	}
	if metadata.PageSize != input.Filters.GetPageSize() {
		t.Error("method get all products don't working well")
	}
	for _, product := range products {
		t.Log(product)
	}
}

func TestUpdateProduct(t *testing.T) {
	product := &proto.Product{
		Id:          id,
		Name:        "Apple",
		Price:       1200,
		Description: "Apple from Almaty city",
		Category:    "Fruit",
		Quantity:    3,
	}
	err := products.Update(product)
	if err != nil {
		t.Fatalf("error acquired while updating product. %s", err.Error())
	}
	result, _ := products.Get(product.Id)
	if result.Name != product.Name && result.Price != product.Price {
		t.Error("returned another product")
	}
}

func TestDeleteProduct(t *testing.T) {
	err := products.Delete(id)
	if err != nil {
		t.Fatalf("error acquired while deleting product. %s", err.Error())
	}
	_, err = products.Get(id)
	if !errors.Is(err, ErrRecordNotFound) {
		t.Error("product should be deleted, but it is not")
	}
}

func getEnvironmentVar(key string) string {
	godotenv.Load("..\\..\\.env")
	return os.Getenv(key)
}

func loadTestConfiguration() *config.Config {
	port, _ := strconv.Atoi(getEnvironmentVar("PORT"))
	return &config.Config{
		Port: port,
		Env:  "development",
		DB: struct {
			DSN          string
			MaxOpenConns int
			MaxIdleConns int
			MaxIdleTime  string
		}{
			DSN:          getEnvironmentVar("DB_DSN"),
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxIdleTime:  "15m",
		},
	}
}
