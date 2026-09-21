package task

import (
	"time"
)

type Task struct {
	ID int `json:"id"`
	Description string `json:"description"`
	Status string `json:"status"`// todo, in-progress, done
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

