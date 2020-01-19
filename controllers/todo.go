package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	orm "github.com/go-pg/pg/v9/orm"
	"github.com/go-pg/pg/v9"
	guuid "github.com/google/uuid"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Completed string      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Create User Table
func CreateTodoTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	createError := db.CreateTable(&Todo{}, opts)
	if createError != nil {
		log.Printf("Error while creating todo table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Todo table created")
	return nil
}

// INITIALIZE DB CONNECTION (TO AVOID TOO MANY CONNECTION)
var dbConnect *pg.DB
func InitiateDB(db *pg.DB) {
	dbConnect = db
}

func GetAllTodos(c *gin.Context) {
	var todos []Todo
	err := dbConnect.Model(&todos).Select()

	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Todos",
		"data": todos,
	})
	return
}

func CreateTodo(c *gin.Context) {
	var todo Todo
	c.BindJSON(&todo)
	title := todo.Title
	body := todo.Body
	completed := todo.Completed
	id := guuid.New().String()

	insertError := dbConnect.Insert(&Todo{
		ID:         id,
		Title:       title,
		Body:      body,
		Completed:     completed,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	if insertError != nil {
		log.Printf("Error while inserting new todo into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Todo created Successfully",
	})
	return
}

func GetSingleTodo(c *gin.Context) {
	todoId := c.Param("todoId")
	todo := &Todo{ID: todoId}
	err := dbConnect.Select(todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Todo",
		"data": todo,
	})
	return
}

func EditTodo(c *gin.Context) {
	todoId := c.Param("todoId")
	var todo Todo
	c.BindJSON(&todo)
	completed := todo.Completed

	_, err := dbConnect.Model(&Todo{}).Set("completed = ?", completed).Where("id = ?", todoId).Update()
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"message":  "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Todo Edited Successfully",
	})
	return
}

func DeleteTodo(c *gin.Context) {
	todoId := c.Param("todoId")
	todo := &Todo{ID: todoId}

	err := dbConnect.Delete(todo)
	if err != nil {
		log.Printf("Error while deleting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Todo deleted successfully",
	})
	return
}
