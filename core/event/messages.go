package event

type MessageSubject string

const (
	NotifyUI         MessageSubject = "ui.notify"
	TaskCreated      MessageSubject = "task.created"
	JobCreated       MessageSubject = "job.created"
	TaskExecutionLog MessageSubject = "task.execution.log"
)

type Envelope struct {
	Subject MessageSubject `json:"subject"`
	Content interface{}    `json:"content"`
}

type UINotificationType string

const (
	UITaskCreated UINotificationType = "task.created"
	UITaskRemoved UINotificationType = "task.removed"
	UIJobExecuted UINotificationType = "job.executed"
	UIJobCreated  UINotificationType = "task.execution.log"
	UIJobRemoved  UINotificationType = "task.execution.log"
)

type UINotification struct {
	Type    UINotificationType `json:"type"`
	Content interface{}        `json:"content"`
}
