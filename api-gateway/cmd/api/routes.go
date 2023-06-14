package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/products", app.addProductHandler)
	router.HandlerFunc(http.MethodGet, "/v1/products", app.listProductsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/products/:id", app.showProductHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/products/:id", app.updateProductHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/products/:id", app.deleteProductHandler)

	return app.recoverPanic(app.rateLimit(router))
}
