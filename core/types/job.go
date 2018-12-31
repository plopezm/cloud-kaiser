package types

import "time"

// JobStatus of the current process
type JobStatus string

const (
	// STOPPED The process is stopped
	RunnableStatusStopped = "STOPPED"
	// RUNNING The process is currently running
	RunnableStatusRunning = "RUNNING"
)

// Job Represents executable job
type Job struct {
	Name       string             `json:"name"`
	Version    string             `json:"version"`
	CreatedAt  time.Time          `json:"created_at"`
	Parameters []JobArgs          `json:"params"`
	Activation JobActivation      `json:"activation"`
	Entrypoint string             `json:"entrypoint"`
	Tasks      map[string]JobTask `json:"tasks"`
	Status     JobStatus          `json:"status"`
}

// JobActivation Represents an activation type
type JobActivation struct {
	Type JobActivationType `json:"type"`
	// Timer represents an ISO 8601 Duration
	Duration string `json:"duration"`
}

// JobActivationType Defines types for launching jobs
type JobActivationType string

const (
	//LOCAL Normal job executed once is received
	LOCAL JobActivationType = "local"
	//REMOTE This job will be executed by request
	REMOTE JobActivationType = "remote"
)

// JobArgs Represents the input arguments to the executor
type JobArgs struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// JobTask Represents a job task to be performed
type JobTask struct {
	Task
	OnSuccess string `json:"onSuccess"`
	OnFailure string `json:"onFailure"`
}
