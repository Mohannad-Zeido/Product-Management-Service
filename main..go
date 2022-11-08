package main

import (
	DAL "ProductManagementService/DBAccess/Database"
	inv "ProductManagementService/DBAccess/Database/InventoryTable"
	invModel "ProductManagementService/DBAccess/Database/InventoryTable/Model"
	pr "ProductManagementService/DBAccess/Database/ProductTable"
	_ "ProductManagementService/DBAccess/Database/ProductTable/Model"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var db *sql.DB

func AddProduct(w http.ResponseWriter, r *http.Request) {

	var product invModel.Inventory

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&product)
	if err != nil {
		log.Fatal(err)
	}
	result := inv.UpsertProduct(product, db)

	_, err = fmt.Fprintf(w, "Welcome to the HomePage! %t ", result)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Endpoint Hit: AddProduct")
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	product := pr.GetProductByBarcode(query.Get("prod"), db)

	responseJSON, err := json.MarshalIndent(product, "", "    ")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(string(responseJSON))
	//w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Fatal(err)
		return
	}

	//fmt.Fprintf(w, "The product you want is %s with barcode %s and %d in stock", product.Name, product.Barcode, product.Stock)
	fmt.Println("Endpoint Hit: product")
}

func handleRequests() {

	db = DAL.ConnectToDatabase()

	http.HandleFunc("/get/product", getProduct)
	http.HandleFunc("/add/product", AddProduct)
	fmt.Println("Server ready and listening on port 1000")
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
