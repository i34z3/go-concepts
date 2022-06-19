package main

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products []Product

func main() {
	e := echo.New()
	e.GET("/products", listProducts)
	e.POST("/product", createProduct)
	e.Logger.Fatal(e.Start(":9191"))
}

func listProducts(c echo.Context) error {
	return c.JSON(200, products)
}

func createProduct(c echo.Context) error {
	product := Product{}
	c.Bind(&product)
	err := persistProduct(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusCreated, product)
}

func persistProduct(product Product) error {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO products (name, price) VALUES ($1, $2)")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}
