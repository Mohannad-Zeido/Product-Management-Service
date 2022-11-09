package main

import (
	DAL "ProductManagementService/DBAccess/Database"
	inv "ProductManagementService/DBAccess/Database/InventoryTable"
	invModel "ProductManagementService/DBAccess/Database/InventoryTable/Model"
	_ "ProductManagementService/DBAccess/Database/ProductTable/Model"
	"ProductManagementService/Services"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	handleRequests()
}

func handleRequests() {

	db = DAL.ConnectToDatabase()
	r := mux.NewRouter()
	r.HandleFunc("/products/{barcode}", getProduct).Methods("GET")
	r.HandleFunc("/inventory/add", addProduct)
	http.Handle("/", r)

	fmt.Println("Server ready and listening on port 1000")
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func addProduct(w http.ResponseWriter, r *http.Request) {

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

	fmt.Println("Endpoint Hit: addProduct")
}

func getProduct(w http.ResponseWriter, r *http.Request) {

	log.Println("Endpoint Hit: product")

	vars := mux.Vars(r)

	productService := Services.ProductService{DbConnection: db}
	retrievedProduce, err := productService.GetProductByBarcode(vars["barcode"])
	if err != nil {
		log.Fatal(err)
		return
	}

	responseJSON, err := json.MarshalIndent(retrievedProduce, "", "    ")
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println(string(responseJSON))

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Fatal(err)
		return
	}
}
