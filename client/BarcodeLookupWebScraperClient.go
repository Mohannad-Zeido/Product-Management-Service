package client

import (
	"io"
	"log"
	"net/http"
)

func GetProduct(barcode string) {
	url := "localhost:1234/lookup?barcode=" + barcode
	resp, getErr := http.Get(url)
	if getErr != nil {
		log.Fatal(getErr)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)
	log.Println(bodyString)
}
