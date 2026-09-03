package dto

// ProductResponse is the HTTP presentation shape for a product.
type ProductResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
