package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Product represents a product with an ID, name, and price
type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
type Payment struct {
	CardNumber    string  `json:"cardNumber"`
	ExpirationDate  string  `json:"expirationDate"`
	CVV string `json:"cvv"`
}

var dummy = Product{ID: "1", Name: "Dummy", Price: 1.99}
var dummyPayment =  Payment{CardNumber:"1231323123", ExpirationDate:"12/12/2020", CVV:"123"} 


// products is a slice of Product representing the data source for the application
var products []Product
var payments []Payment

func main() {
	products = append(products, dummy)
	payments = append(payments, dummyPayment)
	// Initialize the Echo instance
	e := echo.New()
// Add the CORS middleware to allow requests from any origin
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
        c.Response().Header().Set(echo.HeaderAccessControlAllowMethods, "GET, PUT, POST, DELETE")
        c.Response().Header().Set(echo.HeaderAccessControlAllowHeaders, "Content-Type, Authorization")
        return next(c)
    }
})
	// Define the routes for the CRUD operations
	e.GET("/products", getProducts)
	e.GET("/api/payments", getPayments)
	e.GET("/products/:id", getProduct)
	e.POST("/products", createProduct)
	e.PUT("/products/:id", updateProduct)
	e.DELETE("/products/:id", deleteProduct)
	e.POST("/api/payment", createPayment)
	// Start the server
	e.Logger.Fatal(e.Start("localhost:8090"))
}

// getProducts is the handler for the GET /products route, which returns a list of all products
func getProducts(c echo.Context) error {
	return c.JSON(http.StatusOK, products)
}

// getProduct is the handler for the GET /products/{id} route, which returns a single product with the given ID
func getProduct(c echo.Context) error {
	id := c.Param("id")
	for _, product := range products {
		if product.ID == id {
			return c.JSON(http.StatusOK, product)
		}
	}
	// If the product is not found, return a 404 error
	return echo.NewHTTPError(http.StatusNotFound, "Product not found")
}

// getPayments is the handler for the GET /payment route, which returns all payments
func getPayments(c echo.Context) error {
	return c.JSON(http.StatusOK, payments)
}
// createProduct is the handler for the POST /products route, which creates a new product
func createProduct(c echo.Context) error {
	var product Product
	err := c.Bind(&product) // Bind the JSON request body into a Product struct
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	products = append(products, product)
	return c.JSON(http.StatusCreated, product)
}
// createPayment is the handler for the POST /payment route, which creates a new payment
func createPayment(c echo.Context) error {
	var payment Payment
	err := c.Bind(&payment) // Bind the JSON request body into a Product struct
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	payments = append(payments, payment)
	return c.JSON(http.StatusCreated, payment)
}

// updateProduct is the handler for the PUT /products/{id} route, which updates an existing product with the given ID
func updateProduct(c echo.Context) error {
	id := c.Param("id")
	for i, product := range products {
		if product.ID == id {
			err := c.Bind(&product) // Bind the JSON request body into the existing Product struct
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			products[i] = product
			return c.JSON(http.StatusOK, product)
		}
	}
	// If the product is not found, return a 404 error
	return echo.NewHTTPError(http.StatusNotFound, "Product not found")
}

// deleteProduct is the handler for the DELETE /products/{id} route, which deletes an existing product with the given ID
func deleteProduct(c echo.Context) error {
	id := c.Param("id")
	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...) // Remove the existing product
			return c.NoContent(http.StatusNoContent)
		}
	}
	// If the product is not found, return a 404 error
	return echo.NewHTTPError(http.StatusNotFound, "Product not found")
}