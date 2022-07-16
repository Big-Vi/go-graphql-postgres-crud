package main

import (
	"log"
	"net/http"
	"os"
	"context"

	"github.com/big-vi/go-graphql-postgres-crud/internal/pg"
	"github.com/big-vi/go-graphql-postgres-crud/graph"
	"github.com/big-vi/go-graphql-postgres-crud/graph/generated"
	"github.com/big-vi/go-graphql-postgres-crud/graph/model"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
) 

const defaultPort = "8000"

func main() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	toDoService := &postgres.ToDoImpl{}

	toDoService.Initialise()

	config := generated.Config{Resolvers: &graph.Resolver{ToDo: toDoService}}
	config.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (interface{}, error) {
		// TODO: Add user
		// if !getCurrentUser(ctx).HasRole(role) {
		// 	// block calling the next resolver
		// 	return nil, fmt.Errorf("Access denied")
		// }

		// For now let it pass through
		return next(ctx)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
