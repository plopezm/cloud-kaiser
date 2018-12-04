package types

import "time"

// Task the simplest executable node
type Task struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	Script    *string   `json:"script"`
}
