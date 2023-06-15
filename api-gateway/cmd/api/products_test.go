package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListProductsHandler(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(testingApplication.listProductsHandler))

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	//fmt.Println(resp.Body)
	server.Close()
}

func TestAddProductHandler(t *testing.T) {
	var input = `
		{
			"name": "Football",
			"price": 4000,
			"description": "description4",
			"category": "Sports",
			"quantity": 0
		}`

	server := httptest.NewServer(http.HandlerFunc(testingApplication.addProductHandler))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", strings.NewReader(input))
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}
