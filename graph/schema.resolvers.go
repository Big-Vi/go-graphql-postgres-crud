package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/big-vi/go-graphql-postgres-crud/graph/generated"
	"github.com/big-vi/go-graphql-postgres-crud/graph/model"
	"github.com/big-vi/go-graphql-postgres-crud/internal/todo"
	"github.com/google/uuid"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	id, err := r.ToDo.Create(input.Text, false)
	if err != nil {
		return nil, err
	}

	todo := &model.Todo{
		ID:   *id,
		Text: input.Text,
		Done: false,
	}
	return todo, nil
}

// UpdateTodo is the resolver for the updateTodo field.
func (r *mutationResolver) UpdateTodo(ctx context.Context, id string, updatedTodo model.NewTodo) (*model.Todo, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid identifier: %s", id)
	}
	err = r.ToDo.Update(id, updatedTodo.Text, *updatedTodo.Done)
	if err != nil {
		return nil, err
	} else {
		log.Printf("To Do with identifier: %s updated", id)
		return &model.Todo{
			ID:   id,
			Text: updatedTodo.Text,
			Done: *updatedTodo.Done,
		}, nil
	}
}

// DeleteTodo is the resolver for the deleteTodo field.
func (r *mutationResolver) DeleteTodo(ctx context.Context, id string) (string, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return "nil", fmt.Errorf("invalid identifier: %s", id)
	}
	itemId, err := r.ToDo.Delete(id)
	if err != nil {
		return "nil", err
	}
	return *itemId, err
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	var items []*model.Todo
	var savedItems []todo.ToDoItem
	savedItems, err := r.ToDo.GetAll()

	if err != nil {
		return nil, err
	}
	for i, savedItem := range savedItems {
		var item model.Todo
		savedItem = savedItems[i]
		item.ID = savedItem.Id
		item.Text = savedItem.Text
		item.Done = savedItem.Done
		items = append(items, &item)
	}
	return items, nil
}

// Todo is the resolver for the todo field.
func (r *queryResolver) Todo(ctx context.Context, id string) (*model.Todo, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid identifier: %s", id)
	}
	item, err := r.ToDo.Get(id)
	if err != nil {
		return nil, err
	} else {
		return &model.Todo{
			ID:   item.Id,
			Text: item.Text,
			Done: item.Done,
		}, nil
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
