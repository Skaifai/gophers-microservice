package main

import "github.com/Skaifai/gophers-microservice/product-service/pkg/proto"

type Server struct {
	proto.UnimplementedProductServiceServer
}

func NewServer() *Server {
	return &Server{}
}
