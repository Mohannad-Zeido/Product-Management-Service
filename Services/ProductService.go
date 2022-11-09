package Services

import (
	"ProductManagementService/DBAccess/Database/ProductTable"
	"ProductManagementService/DBAccess/Database/ProductTable/Model"
	"database/sql"
	"errors"
)

type ProductService struct {
	DbConnection *sql.DB
}

func (ps ProductService) GetProductByBarcode(productBarcode string) (Model.Product, error) {

	retrievedProduct := ProductTable.GetProductByBarcode(productBarcode, ps.DbConnection)

	if retrievedProduct.Barcode != "" {
		return retrievedProduct, nil
	}

	//TODO add call to get product from Barcode Lookup

	return Model.Product{}, errors.New("error getting the product")
}
