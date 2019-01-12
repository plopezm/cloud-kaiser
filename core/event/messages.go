package event

type MessageSubject string

const (
	TaskCreated MessageSubject = "task.created"
)

type Envelope struct {
	Destination MessageSubject
	Content     interface{}
}
