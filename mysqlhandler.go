package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"database/sql"
	"fmt"
)

type myTodo struct {
	id   int
	task string
}

type Todo []myTodo

func main() {
	db, err := sql.Open("mysql", "root:@/go_test")

	if err != nil {
		log.Fatal("Could not connect, error ", err.Error())
		panic(err.Error())
	}
	defer db.Close()

	fmt.Println("Getting All Todos")
	todos := getAll(db)

	for i := 0; i < len(todos); i++ {
		todo := myTodo{id: todos[i].id, task: todos[i].task}
		fmt.Println(todo.task)
	}

	fmt.Println("Get single todo by ID")
	fmt.Println(getById(db, 1))

	fmt.Println("Add new todo")
	AddNewTodo(db, myTodo{id: 0, task: "New Todo"})

	fmt.Println("Delete todo by ID")
	DeleteById(db, 3)
}

// it will return last insertedid which is int64
func AddNewTodo(db *sql.DB, todo myTodo) int64 {
	res, err := db.Exec("INSERT INTO todos (task) VALUES (?)", todo.task)

	if err != nil {
		panic(err.Error())
	}

	ra, _ := res.RowsAffected()
	id, _ := res.LastInsertId()

	fmt.Println("Rows affected", ra, "Last Inserted Id", id)
	return id
}

func DeleteById(db *sql.DB, id int) {
	res, err := db.Exec("DELETE FROM todos where id = ?", id)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(res)
}

func getById(db *sql.DB, id int) (todo myTodo) {

	row := db.QueryRow("SELECT * FROM todos where id = ?", id)

	err := row.Scan(&todo.id, &todo.task)
	if err != nil {
		panic(err.Error())
	}

	return
}

func getAll(db *sql.DB) []myTodo {
	rows, err := db.Query("SELECT * from todos")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	todos := Todo{}

	for rows.Next() {
		todo := myTodo{}
		err = rows.Scan(&todo.id, &todo.task)
		if err != nil {
			panic(err.Error())
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return todos
}
