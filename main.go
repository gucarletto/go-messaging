package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/gucarletto/go-messaging/internal/infra/akafka"
	"github.com/gucarletto/go-messaging/internal/infra/repository"
	"github.com/gucarletto/go-messaging/internal/infra/web"
	"github.com/gucarletto/go-messaging/internal/usecase"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/products")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize repositories
	productRepository := repository.NewProductRepositoryMySQL(db)
	// Initialize use cases
	createProductUseCase := usecase.NewCreateProductUseCase(productRepository)
	listProductsUseCase := usecase.NewListProductsUseCase(productRepository)

	productHandlers := web.NewProductHandler(createProductUseCase, listProductsUseCase)

	r := chi.NewRouter()
	r.Post("/products", productHandlers.CreateProductHandler)
	r.Get("/products", productHandlers.ListProductsHandler)
	// Start HTTP server
	go http.ListenAndServe(":8000", r)

	// Initialize Kafka consumer
	msgChan := make(chan *kafka.Message)
	go akafka.Consume([]string{"products"}, "host.docker.internal:9094", msgChan)

	for msg := range msgChan {
		dto := usecase.CreateProductInputDTO{}
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			fmt.Println("Error on message:", err)
			continue
		}
		_, err = createProductUseCase.Execute(dto)
		if err != nil {
			fmt.Println("Error on use case:", err)
			continue
		}
	}

}
