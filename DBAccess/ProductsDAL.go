package DBAccess

import (
	"ProductManagementService/DBAccess/Database/ProductTable"
	"ProductManagementService/DBAccess/Model"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"log"
)

type DbConfig struct {
	Username string
	Password string
	Port     int
	Database string
}

func ConnectToDatabase() *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", Config.Username, Config.Password, Config.Port, Config.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func UpsertProduct(product Model.Product, db *sql.DB) bool {
	result, err := updateProduct(product, db)
	if err != nil {
		log.Fatalln(err)
		return false
	}

	if result == false {
		return insertProduct(product, db)
	}

	return result
}

func updateProduct(product Model.Product, db *sql.DB) (bool, error) {
	existingProduct := GetProductByBarcode(product.Barcode, db)

	if existingProduct.Name != "" {
		updateQuery, params, err := sq.Update(ProductTable.TableName).
			Set(ProductTable.StockColumn, existingProduct.Stock+product.Stock).
			Where(sq.Eq{ProductTable.BarcodeColumn: existingProduct.Barcode}).
			RunWith(db).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return false, err
		}

		_, err = db.Query(updateQuery, params...)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func insertProduct(product Model.Product, db *sql.DB) bool {
	insertQuery, params, err := sq.Insert(ProductTable.TableName).
		Columns(ProductTable.BarcodeColumn, ProductTable.NameColumn, ProductTable.StockColumn).
		Values(product.Barcode, product.Name, product.Stock).
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
		From(ProductTable.TableName).
		Where(sq.Eq{ProductTable.BarcodeColumn: barcodeToSearch}).
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

	var name string
	var barcode string
	var stock int

	//if row.
	exists := row.Next()
	if !exists {
		return Model.Product{}
	}

	err = row.Scan(&barcode, &name, &stock)
	if err != nil {
		log.Fatalln(err)
		return Model.Product{}
	}

	return Model.Product{
		Barcode: barcode,
		Name:    name,
		Stock:   stock,
	}
}
