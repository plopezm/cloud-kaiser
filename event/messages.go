package event

type MessageAddress string

const (
	TaskCreated MessageAddress = "task.created"
)

type Envelope struct {
	Destination MessageAddress
	Content     interface{}
}