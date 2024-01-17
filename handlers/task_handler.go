// Upwork test task. will be deleted soon...

package handlers

import (
	"encoding/json"
	"net/http"

	"go_psql_chi_task/models"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

func RegisterTaskHandlers(r chi.Router, db *pgxpool.Pool) {
	r.Get("/tasks", GetTasks(db))
	r.Get("/tasks/{id}", GetTask(db))
	r.Post("/tasks", CreateTask(db))
	r.Put("/tasks/{id}", UpdateTask(db))
	r.Delete("/tasks/{id}", DeleteTask(db))
}

func GetTasks(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(r.Context(), "SELECT * FROM tasks")
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer rows.Close()

		tasks := []models.Task{}
		for rows.Next() {
			var task models.Task
			if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDateTime); err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			tasks = append(tasks, task)
		}

		respondWithJSON(w, http.StatusOK, tasks)
	}
}

func GetTask(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var task models.Task
		err := db.QueryRow(r.Context(), "SELECT id, title, description, priority, due_date_time FROM tasks WHERE id = $1", id).Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDateTime)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Task not found")
			return
		}

		respondWithJSON(w, http.StatusOK, task)
	}
}

func CreateTask(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&task); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		err := db.QueryRow(r.Context(), "INSERT INTO tasks(title, description, priority, due_date_time) VALUES($1, $2, $3, $4) RETURNING id", task.Title, task.Description, task.Priority, task.DueDateTime).Scan(&task.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, task)
	}
}

func UpdateTask(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var task models.Task
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&task); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		_, err := db.Exec(r.Context(), "UPDATE tasks SET title=$1, description=$2, priority=$3, due_date_time=$4 WHERE id=$5", task.Title, task.Description, task.Priority, task.DueDateTime, id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Task updated successfully"})
	}
}

func DeleteTask(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		_, err := db.Exec(r.Context(), "DELETE FROM tasks WHERE id=$1", id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
	}
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
