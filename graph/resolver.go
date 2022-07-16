package graph

import "github.com/big-vi/go-graphql-postgres-crud/internal/todo"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	ToDo todo.ToDo
}
