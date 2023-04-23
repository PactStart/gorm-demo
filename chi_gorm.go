package main

import (
	"context"
	"encoding/json"
	"gorm.io/driver/mysql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *gorm.DB

func initDB() {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todo{})
}

func main() {
	initDB()

	r := chi.NewRouter()

	r.Get("/todos", getTodos)
	r.Post("/todos", createTodo)
	r.Route("/todos/{id}", func(r chi.Router) {
		r.Use(todoCtx)
		r.Get("/", getTodo)
		r.Put("/", updateTodo)
		r.Delete("/", deleteTodo)
	})

	http.ListenAndServe(":8080", r)
}

func todoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		todo := Todo{}
		if err := db.First(&todo, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.NotFound(w, r)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "todo", todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo

	if err := db.Find(&todos).Error; err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Failed to get todos: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	if err := db.Create(&todo).Error; err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Failed to create new todo: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value("todo").(Todo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value("todo").(Todo)

	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	todo.Title = updatedTodo.Title
	todo.Completed = updatedTodo.Completed

	if err := db.Save(&todo).Error; err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Failed to update todo: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value("todo").(Todo)

	if err := db.Delete(&todo).Error; err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("Failed to delete todo: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
