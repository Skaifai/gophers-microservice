package main

import (
	"api-gateway/internal/validator"
	productServiceProto "github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"testing"
)

func TestValidateProduct(t *testing.T) {
	v := validator.New()

	noNameProduct := &productServiceProto.Product{
		Name:        "",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	ValidateProduct(v, noNameProduct)

	expected := v.Valid()

	if expected != false {
		t.Errorf("ValidateProduct(noNameProduct) returned unexpected value: got %v, expected %s", expected, "false")
	}
}

func TestTableDrivenValidateProduct(t *testing.T) {
	noNameProduct := &productServiceProto.Product{
		Name:        "",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	longNameProduct := &productServiceProto.Product{
		Name:        "This A Very Long Name That Has More Than Twenty Bytes In It",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	badPriceProduct := &productServiceProto.Product{
		Name:        "GoodName",
		Price:       -100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	emptyDescriptionProduct := &productServiceProto.Product{
		Name:        "GoodName",
		Price:       100.0,
		Description: "",
		Category:    "Category",
		Quantity:    5,
	}

	emptyCategoryProduct := &productServiceProto.Product{
		Name:        "GoodName",
		Price:       100.0,
		Description: "Descrpition",
		Category:    "",
		Quantity:    5,
	}

	badQuantityProduct := &productServiceProto.Product{
		Name:        "GoodName",
		Price:       100.0,
		Description: "",
		Category:    "Category",
		Quantity:    -1,
	}

	perfectProduct := &productServiceProto.Product{
		Name:        "GoodName",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	var tests = []struct {
		name     string
		input    *productServiceProto.Product
		expected bool
	}{
		{
			"ValidateProduct(noNameProduct) must return false",
			noNameProduct,
			false,
		},
		{
			"ValidateProduct(longNameProduct) must return false",
			longNameProduct,
			false,
		},
		{
			"ValidateProduct(badPriceProduct) must return false",
			badPriceProduct,
			false,
		},
		{
			"ValidateProduct(emptyDescriptionProduct) must return false",
			emptyDescriptionProduct,
			false,
		},
		{
			"ValidateProduct(emptyCategoryProduct) must return false",
			emptyCategoryProduct,
			false,
		},
		{
			"ValidateProduct(badQuantityProduct) must return false",
			badQuantityProduct,
			false,
		},
		{
			"ValidateProduct(perfectProduct) must return true",
			perfectProduct,
			true,
		},
	}

	for _, tst := range tests {
		v := validator.New()
		t.Run(tst.name, func(t *testing.T) {
			ValidateProduct(v, tst.input)
			result := v.Valid()
			if result != tst.expected {
				t.Errorf("Expected %v got %v", tst.expected, result)
			}
		})
	}
}

func TestSetCurrentStatus(t *testing.T) {
	someProduct := &productServiceProto.Product{
		Name:        "Hello",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	SetStatus(someProduct.Quantity, someProduct)

	expected := someProduct.IsAvailable

	if expected != true {
		t.Errorf("SetStatus(noNameProduct.Quantity) returned unexpected value: got %v, expected %s", expected, "true")
	}
}

func TestSetNewStatus(t *testing.T) {
	someProduct := &productServiceProto.Product{
		Name:        "Product",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	SetStatus(0, someProduct)

	expected := someProduct.IsAvailable

	if expected != false {
		t.Errorf("SetStatus(noNameProduct.Quantity) returned unexpected value: got %v, expected %s", expected, "true")
	}
}
