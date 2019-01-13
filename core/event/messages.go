package event

type MessageSubject string

const (
	TaskCreated      MessageSubject = "task.created"
	JobCreated       MessageSubject = "job.created"
	TaskExecutionLog MessageSubject = "task.execution.log"
)

type Envelope struct {
	Subject MessageSubject `json:"subject"`
	Content interface{}    `json:"content"`
}
