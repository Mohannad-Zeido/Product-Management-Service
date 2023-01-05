package ProductTable

import (
	"ProductManagementService/DBAccess/Database/ProductTable/Model"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"log"
)

func InsertProduct(product Model.Product, db *sql.DB) bool {

	insertQuery, params, err := sq.Insert(TableName).
		Columns(BarcodeColumn, NameColumn, ImageUrlColumn, CategoryColumn).
		Values(product.Barcode, product.Name, product.ImageUrl, product.Category, product.Quantity).
		RunWith(db).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Fatalln(err)
		return false
	}

	_, err = db.Query(insertQuery, params...)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	return true
}

func GetProductByBarcode(barcodeToSearch string, db *sql.DB) Model.Product {

	seql, params, err := sq.Select("*").
		From(TableName).
		Where(sq.Eq{BarcodeColumn: barcodeToSearch}).
		RunWith(db).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		log.Fatalln(err)
	}

	row, err := db.Query(seql, params...)
	if err != nil {
		log.Fatalln(err)
	}
	defer row.Close()

	var returnedProduct = Model.Product{}

	exists := row.Next()
	if !exists {
		return Model.Product{}
	}

	err = row.Scan(&returnedProduct.Barcode, &returnedProduct.Name, &returnedProduct.ImageUrl, &returnedProduct.Category, &returnedProduct.Quantity)
	if err != nil {
		log.Fatalln(err)
		return Model.Product{}
	}

	return returnedProduct
}
