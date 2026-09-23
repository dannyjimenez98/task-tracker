package task

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      string     `json:"status"` // todo, in-progress, done
	CreatedAt   CustomTime `json:"createdAt"`
	UpdatedAt   CustomTime `json:"updatedAt"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type CustomTime struct {
	time.Time
}

// MarshalJSON is a custom implementation for the defined CustomTime type.
// Time is still written to JSON in RFC3339, but without nanoseconds.
// Timestamps of type CustomTime displayed in JSON as YYYY-MM-DDTHH:MM:SS±HH:MM
func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Time.Truncate(time.Nanosecond).Format(time.RFC3339))
}
