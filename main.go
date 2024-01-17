// main.go: boots up the task as explained in the doc
// Upwork test task. will be deleted soon...

package main

import (
	"net/http"

	"go_psql_chi_task/database"
	"go_psql_chi_task/handlers"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	db := database.InitDB()
	defer db.Close()

	r := chi.NewRouter()
	r.Get("/", handlers.ServeIndex)
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	handlers.RegisterTaskHandlers(r, db)

	http.ListenAndServe(":8080", r)
}
