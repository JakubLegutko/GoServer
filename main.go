package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Product represents a product with an ID, name, and price
type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}
type Payment struct {
	CardNumber     string `json:"cardNumber"`
	ExpirationDate string `json:"expirationDate"`
	CVV            string `json:"cvv"`
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "jon" || password != "shhh!" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		"Jon Snow",
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

var dummy = Product{ID: "1", Name: "Dummy", Price: 1.99}
var dummyPayment = Payment{CardNumber: "1231323123", ExpirationDate: "12/12/2020", CVV: "123"}

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
	e.GET("/", helloWorld)
	e.GET("/api/payments", getPayments)
	e.GET("/products/:id", getProduct)
	e.POST("/products", createProduct)
	e.PUT("/products/:id", updateProduct)
	e.DELETE("/products/:id", deleteProduct)
	e.POST("/api/payment", createPayment)
	e.POST("/login", login)
	e.GET("/accessible", accessible)
	// Restricted group
	r := e.Group("/restricted")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("", restricted)
	// Start the server
	e.Logger.Fatal(e.Start("0.0.0.0:8090"))
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
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

// getPayments is the handler for the GET /api/payment route, which returns all payments
func getPayments(c echo.Context) error {
	return c.JSON(http.StatusOK, payments)
}

// createProduct is the handler for the POST /products route, which creates a new product
func createProduct(c echo.Context) error {
	var product Product
	// Check that the request isn't empty
	if c.Request().ContentLength == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Request body is empty")
	}
	err := c.Bind(&product) // Bind the JSON request body into a Product struct
	if product.ID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing product ID")
	}
	if product.Price == 0 {
		return echo.NewHTTPError(402, "Missing product price!")
	}
	for i := 0; i < len(products); i++ {
		if product.ID == products[i].ID {
			return echo.NewHTTPError(409, "Product ID already exists")
		}
	}
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
	if c.Request().ContentLength == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Payment body is empty")
	}
	if payment.CVV == "" {
		return echo.NewHTTPError(402, "Missing CVV")
	}
	if len(payment.CVV) != 3 {
		return echo.NewHTTPError(409, "Invalid CVV")
	}
	if len(payment.CardNumber) > 16 {
		return echo.NewHTTPError(409, "Invalid Card Number!")
	}

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
