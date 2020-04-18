package types

import "time"

type TaskExecutionLog struct {
	JobName     string    `json:"jobName"`
	JobVersion  string    `json:"jobVersion"`
	TaskName    string    `json:"taskName"`
	TaskVersion string    `json:"taskVersion"`
	Line        string    `json:"line"`
	Ts          time.Time `json:"ts"`
}
