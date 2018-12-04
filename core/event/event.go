package event

import "github.com/plopezm/cloud-kaiser/core/types"

type EventStore interface {
	Close()
	PublishTaskCreated(types.Task) error
	OnTaskCreated(MessageAddress, func(task types.Task)) error
}

var impl EventStore

func SetEventStore(es EventStore) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishTaskCreated(task types.Task) error {
	return impl.PublishTaskCreated(task)
}

func OnTaskCreated(ma MessageAddress, f func(task types.Task)) error {
	return impl.OnTaskCreated(ma, f)
}
