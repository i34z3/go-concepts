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
	e.GET("/products", getProducts)
	e.POST("/product", postProduct)
	e.Logger.Fatal(e.Start(":9191"))
}

func getProducts(c echo.Context) error {
	products = readProducts()
	return c.JSON(http.StatusOK, products)
}

func readProducts() []Product {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	smtp, err := db.Prepare("SELECT name, price FROM products")
	if err != nil {
		panic(err)
	}
	rows, err := smtp.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var price float64
		err = rows.Scan(&name, &price)
		if err != nil {
			panic(err)
		}
		products = append(products, Product{Name: name, Price: price})
	}
	return products
}

func postProduct(c echo.Context) error {
	product := Product{}
	c.Bind(&product)
	err := createProduct(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusCreated, product)
}

func createProduct(product Product) error {
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
