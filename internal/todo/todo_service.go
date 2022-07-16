package todo

import "time"

type ToDoItem struct {
	Id        string
	Text      string
	Done      bool
	CreatedOn time.Time
	UpdatedOn *time.Time
}

type ToDo interface {
	Initialise() error
	Create(Text string, Done bool) (*string, error)
	Update(id string, text string, Done bool) error
	Get(id string) (*ToDoItem, error)
	GetAll() ([]ToDoItem, error)
	Delete(id string) (*string, error)
}