package server

import (
	"context"
	"fmt"
	"github.com/Skaifai/gophers-microservice/product-service/cmd/utils"
	"github.com/Skaifai/gophers-microservice/product-service/config"
	"github.com/Skaifai/gophers-microservice/product-service/internal/logger"
	"github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"testing"
)

var server = func() *Server {
	cfg := loadTestConfiguration()
	db, err := utils.OpenDB(cfg)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	publisher, err := logger.NewPublisher()
	if err != nil {
		log.Fatalf("failed to create publisher: %v", err)
	}
	return NewServer(db, publisher)
}()

func TestServer_ShowProduct(t *testing.T) {
	req := &proto.ShowProductRequest{
		Id: 4,
	}
	res, err := server.ShowProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("error acquired while showing product. %s", err.Error())
	}
	if res.GetProduct() == nil {
		t.Fatalf("product is nil. %s", err.Error())
	}
}

func TestServer_ListProducts(t *testing.T) {
	req := &proto.ListProductsRequest{
		Name:     "",
		Category: "",
		Filters: &proto.Filters{
			Page:     1,
			PageSize: 20,
			Sort:     "id",
			SortSafeList: []string{"id", "name", "category", "price", "is_available", "creation_date",
				"-id", "-name", "-category", "-price", "-is_available", "-creation_date"},
		},
	}
	res, err := server.ListProducts(context.Background(), req)
	if err != nil {
		t.Fatalf("error acquired while listing product. %s", err.Error())
	}
	if res.GetProducts() == nil {
		t.Fatalf("products is nil. %s", err.Error())
	}
}

func TestServer_UpdateProduct(t *testing.T) {
	req := &proto.UpdateProductRequest{
		Product: &proto.Product{
			Id:          4,
			Name:        "Apple",
			Price:       1200,
			Description: "Apple from Almaty city",
			Category:    "Fruit",
			Quantity:    3,
		},
	}
	res, err := server.UpdateProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("error acquired while updating product. %s", err.Error())
	}
	expected := fmt.Sprintf("Product has been successfully updated with id: %d", req.GetProduct().GetId())
	if res.GetMessage() != expected {
		t.Errorf("error acquired while updating product. %s", err.Error())
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
