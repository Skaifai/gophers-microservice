package main

import (
	"context"
	"fmt"
	"github.com/Skaifai/gophers-microservice/api-gateway/internal/validator"
	productServiceProto "github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func (app *application) addProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string  `json:"name"`
		Price       float32 `json:"price"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Quantity    int32   `json:"quantity"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product := &productServiceProto.Product{
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
		Category:    input.Category,
		Quantity:    input.Quantity,
	}
	SetStatus(product.Quantity, product)

	v := validator.New()
	if ValidateProduct(v, product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := app.productServiceClient.AddProduct(ctx, &productServiceProto.AddProductRequest{
		Product: product,
	})
	if err != nil {
		errorStatus, _ := status.FromError(err)
		switch {
		case errorStatus.Code() == codes.DeadlineExceeded:
			app.deadlineExceededResponse(w, r, err)
		case errorStatus.Code() == codes.Unavailable:
			app.serviceUnavailableResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"product": response.Product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	return
}

func (app *application) showProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product, err := app.productServiceClient.ShowProduct(ctx, &productServiceProto.ShowProductRequest{
		Id: id,
	})
	if err != nil {
		errorStatus, _ := status.FromError(err)
		switch {
		case errorStatus.Code() == codes.DeadlineExceeded:
			app.deadlineExceededResponse(w, r, err)
		case errorStatus.Code() == codes.Unavailable:
			app.serviceUnavailableResponse(w, r, err)
		case errorStatus.Code() == codes.NotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Encode the struct to JSON and send it as the HTTP response.
	// using envelope
	err = app.writeJSON(w, http.StatusOK, product, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listProductsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Category string
		productServiceProto.Filters
	}
	// v := validator.New()
	qs := r.URL.Query()
	input.Name = app.readString(qs, "name", "")
	fmt.Println(input.Name)
	input.Category = app.readString(qs, "category", "")
	input.Filters.Page = int32(app.readInt(qs, "page", 1))
	input.Filters.PageSize = int32(app.readInt(qs, "page_size", 20))
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "name", "category", "price", "is_available", "creation_date",
		"-id", "-name", "-category", "-price", "-is_available", "-creation_date"}

	// if data.ValidateFilters(v, input.Filters); !v.Valid() {
	// 	app.failedValidationResponse(w, r, v.Errors)
	// 	return
	// }
	// Call the GetAll() method to retrieve the movies, passing in the various filter
	// parameters.
	// Accept the metadata struct as a return value.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := app.productServiceClient.ListProducts(ctx, &productServiceProto.ListProductsRequest{
		Name:     input.Name,
		Category: input.Category,
		Filters:  &input.Filters,
	})
	if err != nil {
		errorStatus, _ := status.FromError(err)
		switch {
		case errorStatus.Code() == codes.DeadlineExceeded:
			app.deadlineExceededResponse(w, r, err)
		case errorStatus.Code() == codes.Unavailable:
			app.serviceUnavailableResponse(w, r, err)
		case errorStatus.Code() == codes.NotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Include the metadata in the response envelope.
	err = app.writeJSON(w, http.StatusOK, envelope{"products": response.Products, "metadata": response.Metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//if productServiceConnection.GetState() != connectivity.Ready {
	//	app.logger.PrintInfo("Failed to update product due to no connection to the product service", map[string]string{
	//		"method:": "updateProductHandler",
	//	})
	//	app.errorResponse(w, r, http.StatusInternalServerError, "Failed to update product due to no connection to the product service")
	//	return
	//}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	productFromDB, err := app.productServiceClient.ShowProduct(ctx, &productServiceProto.ShowProductRequest{
		Id: id,
	})
	if err != nil {
		errorStatus, _ := status.FromError(err)
		switch {
		case errorStatus.Code() == codes.DeadlineExceeded:
			app.deadlineExceededResponse(w, r, err)
		case errorStatus.Code() == codes.Unavailable:
			app.serviceUnavailableResponse(w, r, err)
		case errorStatus.Code() == codes.NotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	product := productFromDB.Product

	var input struct {
		Name        *string  `json:"name"`
		Price       *float32 `json:"price"`
		Description *string  `json:"description"`
		Category    *string  `json:"category"`
		Quantity    *int32   `json:"quantity"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		product.Name = *input.Name
	}

	if input.Price != nil {
		product.Price = *input.Price
	}

	if input.Description != nil {
		product.Description = *input.Description
	}

	if input.Category != nil {
		product.Category = *input.Category
	}

	if input.Quantity != nil {
		product.Quantity = *input.Quantity
		SetStatus(product.Quantity, product)
	}

	v := validator.New()
	if ValidateProduct(v, product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	response, err := app.productServiceClient.UpdateProduct(ctx, &productServiceProto.UpdateProductRequest{
		Product: product,
	})
	if err != nil {
		errorStatus, _ := status.FromError(err)
		switch {
		case errorStatus.Code() == codes.DeadlineExceeded:
			app.deadlineExceededResponse(w, r, err)
		case errorStatus.Code() == codes.Unavailable:
			app.serviceUnavailableResponse(w, r, err)
		case errorStatus.Code() == codes.NotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.logger.PrintInfo(response.GetMessage(), map[string]string{
		"method": "updateProductHandler",
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"message": response.GetMessage(), "product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//if productServiceConnection.GetState() != connectivity.Ready {
	//	app.logger.PrintInfo("Failed to delete product due to no connection to the product service", map[string]string{
	//		"method:": "deleteProductHandler",
	//	})
	//	app.errorResponse(w, r, http.StatusInternalServerError, "Failed to delete product due to no connection to the product service")
	//	return
	//}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := app.productServiceClient.DeleteProduct(ctx, &productServiceProto.DeleteProductRequest{
		Id: id,
	})
	if err != nil {
		errorStatus, _ := status.FromError(err)
		switch {
		case errorStatus.Code() == codes.DeadlineExceeded:
			app.deadlineExceededResponse(w, r, err)
		case errorStatus.Code() == codes.Unavailable:
			app.serviceUnavailableResponse(w, r, err)
		case errorStatus.Code() == codes.NotFound:
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.logger.PrintInfo(response.GetMessage(), map[string]string{
		"method": "deleteProductHandler",
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"message": response.GetMessage()}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
