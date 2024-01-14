// Upwork test task. will be deleted soon...

package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	DueDateTime time.Time `json:"due_date_time"`
}

var DB *pgxpool.Pool

func main() {
	initDB()
	defer DB.Close()
	r := chi.NewRouter()
	r.Get("/", serveIndex)
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/tasks", GetTasks)
	r.Get("/tasks/{id}", GetTask)
	r.Post("/tasks", CreateTask)
	r.Put("/tasks/{id}", UpdateTask)
	r.Delete("/tasks/{id}", DeleteTask)

	http.ListenAndServe(":8080", r)
}

func initDB() {
	connectionString := "postgresql://postgres:admin@localhost/postgres?sslmode=disable"

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	DB, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(r.Context(), "SELECT * FROM tasks")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDateTime); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		tasks = append(tasks, task)
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var task Task
	err := DB.QueryRow(r.Context(), "SELECT id, title, description, priority, due_date_time FROM tasks WHERE id = $1", id).Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDateTime)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err := DB.QueryRow(r.Context(), "INSERT INTO tasks(title, description, priority, due_date_time) VALUES($1, $2, $3, $4) RETURNING id", task.Title, task.Description, task.Priority, task.DueDateTime).Scan(&task.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	_, err := DB.Exec(r.Context(), "UPDATE tasks SET title=$1, description=$2, priority=$3, due_date_time=$4 WHERE id=$5", task.Title, task.Description, task.Priority, task.DueDateTime, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Task updated successfully"})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := DB.Exec(r.Context(), "DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
