package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Active       bool
	Role         string
	Password     string
	Birthday     time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func SetDBMiddleware(next http.Handler) http.Handler {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	db.Debug()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeoutContext, _ := context.WithTimeout(context.Background(), time.Second)
		ctx := context.WithValue(r.Context(), "DB", db.WithContext(timeoutContext))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	r := chi.NewRouter()
	r.Use(SetDBMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		db, ok := r.Context().Value("DB").(*gorm.DB)
		if !ok {
			fmt.Println("db not ok")
			return
		}
		var users []User
		db.Find(&users)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
		db, ok := r.Context().Value("DB").(*gorm.DB)
		if !ok {
			return
		}
		var user User
		db.First(&user)
		fmt.Println(user)
		// lots of db operations
	})
	http.ListenAndServe(":8080", r)
}
