// Upwork test task. will be deleted soon...

package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	DueDateTime time.Time `json:"due_date_time"`
}
