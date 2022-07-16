package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib" 
	"github.com/big-vi/go-graphql-postgres-crud/internal/todo"
	"log"
	"os"
)

type ToDoImpl struct {
	DbUserName string
	DbPassword string
	DbURL      string
	DbName     string
}

func (t *ToDoImpl) Initialise() error {
	targetSchemaVersion := 2
	connString := t.getDBConnectionString()
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	
	if err != nil {
		fmt.Print(err)
	}
	version, dirty, err := driver.Version()
	if dirty {
		log.Fatalf("ERROR: The current database schema is reported as being dirty. A manual resolution is needed")
	}
	log.Printf("Target database schema version is: %d and current database schema version is: %d", targetSchemaVersion, version)
	if version != targetSchemaVersion {
		log.Printf("Migrating database schema from version: %d to version %d", version, targetSchemaVersion)
		m, err := migrate.NewWithDatabaseInstance("file://migrations", t.DbName, driver) 
		if err != nil {
			fmt.Print(err)
		}
		err = m.Steps(targetSchemaVersion)
		if err != nil {
			return err
		}
		return nil
	} else {
		log.Println("No database schema migrations need to be performed.")
	}
	if err != nil {
		log.Fatalf("ERROR: Could not determine the current database schema version")
	}
	return nil
}

func (t *ToDoImpl) Create(Text string, Done bool) (*string, error) {
	insertSQL := "insert into todo_item(id, text, done) values ($1, $2, $3)"
	ctx := context.Background()
	dbPool := t.getConnection()
	defer dbPool.Close()
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	id := uuid.New()
	idStr := id.String()
	_, err = tx.Exec(ctx, insertSQL, idStr, Text, Done)
	if err != nil {
		log.Println("ERROR: Could not save the To Do item due to the error:", err)
		rollbackErr := tx.Rollback(ctx)
		log.Fatal("ERROR: Transaction rollback failed due to the error: ", rollbackErr)
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &idStr, nil
}

func (t *ToDoImpl) Update(id string, Text string, Done bool) error {
	updateSQL := "update todo_item set text = $1, done = $2 where id = $3"
	ctx := context.Background()
	dbPool := t.getConnection()
	defer dbPool.Close()
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, updateSQL, Text, Done, id)
	if err != nil {
		log.Println("ERROR: Could not save the To Do item due to the error:", err)
		rollbackErr := tx.Rollback(ctx)
		log.Fatal("ERROR: Transaction rollback failed due to the error: ", rollbackErr)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *ToDoImpl) Get(id string) (*todo.ToDoItem, error) {
	selectSQL := "select id, text, done, creation_timestamp, update_timestamp from todo_item where id = $1"
	dbPool := t.getConnection()
	defer dbPool.Close()
	var todoItem todo.ToDoItem
	err := dbPool.QueryRow(context.Background(), selectSQL, id).Scan(&todoItem.Id, &todoItem.Text, &todoItem.Done, &todoItem.CreatedOn, &todoItem.UpdatedOn)
	if err != nil {
		return nil, err
	}
	return &todoItem, nil
}

func (t *ToDoImpl) Delete(id string) (*string, error) {
	deleteSQL := "delete from todo_item where id = $1"
	ctx := context.Background()
	dbPool := t.getConnection()
	defer dbPool.Close()
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, deleteSQL, id)
	if err != nil {
		log.Println("ERROR: Could not delete the To Do item due to the error:", err)
		rollbackErr := tx.Rollback(ctx)
		log.Fatal("ERROR: Transaction rollback failed due to the error: ", rollbackErr)
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &id, nil

}

func (t *ToDoImpl) GetAll() ([]todo.ToDoItem, error) {
	selectSQL := "select id, text, done, creation_timestamp, update_timestamp from todo_item"
	dbPool := t.getConnection()
	defer dbPool.Close()
	var todoItems []todo.ToDoItem
	rows, err := dbPool.Query(context.Background(), selectSQL)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var todoItem todo.ToDoItem
		err = rows.Scan(&todoItem.Id, &todoItem.Text, &todoItem.Done, &todoItem.CreatedOn, &todoItem.UpdatedOn)
		if err != nil {
			return nil, err
		}
		todoItems = append(todoItems, todoItem)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return todoItems, nil
}

func (t *ToDoImpl) getConnection() *pgxpool.Pool {
	dbPool, err := pgxpool.Connect(context.Background(), t.getDBConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	return dbPool
}

func (t *ToDoImpl) getDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"),
	)
}