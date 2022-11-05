package main

import (
	pr "ProductManagementService/DBAccess"
	"ProductManagementService/DBAccess/Model"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var db *sql.DB

func AddProduct(w http.ResponseWriter, r *http.Request) {

	var product Model.Product

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&product)
	if err != nil {
		log.Fatal(err)
	}
	result := pr.UpsertProduct(product, db)

	fmt.Fprintf(w, "Welcome to the HomePage! %t ", result)
	fmt.Println("Endpoint Hit: AddProduct")
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	product := pr.GetProductByBarcode(query.Get("prod"), db)

	responseJSON, err := json.MarshalIndent(product, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseJSON))
	//w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		return
	}

	//fmt.Fprintf(w, "The product you want is %s with barcode %s and %d in stock", product.Name, product.Barcode, product.Stock)
	fmt.Println("Endpoint Hit: product")
}

func handleRequests() {

	db = pr.ConnectToDatabase()

	http.HandleFunc("/get/product", getProduct)
	http.HandleFunc("/add/product", AddProduct)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	//handleRequests()
	//db := pr.ConnectToDatabase()

	db = pr.ConnectToDatabase()
	pr.UpsertProduct(Model.Product{
		Barcode: "135789",
		Name:    "bro",
		Stock:   1,
	}, db)
	//pr.GetProductByBarcode("135789", db)

	//defer db.Close()

}
