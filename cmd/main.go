package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"service-for-working-with-posts-on-graphql/graph"
	"service-for-working-with-posts-on-graphql/internal/config"
	"service-for-working-with-posts-on-graphql/internal/repositories"
	"service-for-working-with-posts-on-graphql/internal/repositories/memorydb"
	"service-for-working-with-posts-on-graphql/pkq/postgres"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading env file")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:")
	}

	var repository graph.Repository

	if cfg.Env == "pql" {
		postgresClient, err := postgres.NewPostgresClient(cfg)
		if err != nil {
			fmt.Errorf("Failed to create PostgreSQL client")
		}
		repository = repositories.NewPgRepository(postgresClient)
	} else {
		repository = memorydb.NewMemoryRepository()
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Repo: repository}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
