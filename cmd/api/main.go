package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-api/internal/httpapi"
	"go-api/internal/product"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	productStore := product.NewStore(pool)

	mux := http.NewServeMux()
	httpapi.RegisterProductRoutes(mux, productStore)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}