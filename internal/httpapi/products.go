package httpapi

import (
	"encoding/json"
	"net/http"
	"go-api/internal/product"
)

func RegisterProductRoutes(mux *http.ServeMux, store *product.Store) {
	mux.HandleFunc("GET /products", ListProducts(store))
}

func ListProducts(store *product.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := store.List(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}