package main

import (
	"fmt"
	"github.com/h-abranches-dev/daily-expenses-be/service-layer/handlers"
	"net/http"
)

func main() {
	hmux := http.NewServeMux()
	hmux.HandleFunc("/transactions", handlers.TransactionHandlerFunc)
	hmux.HandleFunc("/transactions/", handlers.UpdateTransactionHandlerFunc)
	hmux.HandleFunc("/transactions/types", handlers.TransactionsTypesHandlerFunc)
	hmux.HandleFunc("/entities", handlers.TransactionsEntitiesHandlerFunc)
	hmux.HandleFunc("/categories", handlers.CategoriesHandlerFunc)
	hmux.HandleFunc("/categories/", handlers.UpdateCategoryHandlerFunc)
	hmux.HandleFunc("/transactions/categories/", handlers.TransactionsCategoriesHandlerFunc)

	api := http.Server{
		Addr:    ":8080",
		Handler: hmux,
	}
	fmt.Printf("Listening on port %q\n", api.Addr)
	if err := api.ListenAndServe(); err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}
}
