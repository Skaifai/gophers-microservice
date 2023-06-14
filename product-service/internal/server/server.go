package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Skaifai/gophers-microservice/product-service/internal/data"
	"github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedProductServiceServer
	data.Models
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		Models: data.NewModels(db),
	}
}

func (s *Server) ShowProduct(ctx context.Context, req *proto.ShowProductRequest) (*proto.ShowProductResponse, error) {
	product, err := s.Products.Get(req.GetId())
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "Product not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "Failed to retrieve product: %v", err)
	}

	return &proto.ShowProductResponse{
		Product: product,
	}, nil
}

func (s *Server) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
	products, metadata, err := s.Products.GetAll(req.GetName(), req.GetCategory(), req.GetFilters())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get products: %v", err)
	}

	return &proto.ListProductsResponse{
		Metadata: metadata,
		Products: products,
	}, nil
}

func (s *Server) AddProduct(ctx context.Context, req *proto.AddProductRequest) (*proto.AddProductResponse, error) {
	product := req.GetProduct()
	response, err := s.Products.Insert(product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to add product: %v", err)
	}

	return &proto.AddProductResponse{
		Product: response,
	}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	product := req.GetProduct()

	err := s.Products.Update(product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update product: %v", err)
	}

	return &proto.UpdateProductResponse{
		Message: fmt.Sprintf("Product has been successfully updated with id: %d", product.GetId()),
	}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	err := s.Products.Delete(req.GetId())
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "Product not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "Failed to delete product: %v", err)
	}

	return &proto.DeleteProductResponse{
		Message: fmt.Sprintf("Product has been successfully deleted with id: %d", req.GetId()),
	}, nil
}
