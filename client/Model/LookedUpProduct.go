package Model

type LookedUpProduct struct {
	Barcode  string `json:"barcode"`
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
	Category string `json:"category"`
}
