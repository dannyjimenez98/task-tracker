package task

import "time"


type Task struct {
	ID int
	Description string
	Status string // todo, in-progress, done
	createdAt time.Time
	updatedAt time.Time
}


func NewTask(id int, description string) *Task {
	return &Task{
		ID: id,
		Description: description,
		Status: "todo",
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}
