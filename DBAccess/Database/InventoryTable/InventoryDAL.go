package InventoryTable

import (
	"ProductManagementService/DBAccess/Database/InventoryTable/Model"
	"ProductManagementService/DBAccess/Database/ProductTable"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"log"
)

func UpsertProduct(inventory Model.Inventory, db *sql.DB) bool {
	result, err := updateStock(inventory, db)
	if err != nil {
		log.Fatalln(err)
		return false
	}

	if result == false {
		return insertStock(inventory, db)
	}

	return result
}

func updateStock(inventory Model.Inventory, db *sql.DB) (bool, error) {
	existingProduct := GetStock(inventory.Barcode, db)

	if existingProduct.Barcode != "" {
		updateQuery, params, err := sq.Update(TableName).
			Set(StockColumn, existingProduct.Stock+inventory.Stock).
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

func insertStock(inventory Model.Inventory, db *sql.DB) bool {
	insertQuery, params, err := sq.Insert(TableName).
		Columns(BarcodeColumn, StockColumn).
		Values(inventory.Barcode, inventory.Stock).
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

func GetStock(barcodeToSearch string, db *sql.DB) Model.Inventory {

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

	var inventory = Model.Inventory{}

	exists := row.Next()
	if !exists {
		return Model.Inventory{}
	}

	err = row.Scan(&inventory.Barcode, &inventory.Stock)
	if err != nil {
		log.Fatalln(err)
		return Model.Inventory{}
	}

	return inventory

}
