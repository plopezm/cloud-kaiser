package event

import "github.com/plopezm/cloud-kaiser/types"

type EventStore interface {
	Close()
	PublishTaskCreated(types.JobTask) error
	OnTaskCreated(MessageAddress, func(task types.JobTask)) error
}

var impl EventStore

func SetEventStore(es EventStore) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishTaskCreated(task types.JobTask) error {
	return impl.PublishTaskCreated(task)
}

func OnTaskCreated(ma MessageAddress, f func(task types.JobTask)) error {
	return impl.OnTaskCreated(ma, f)
}
