package event

type EventStore interface {
	Close()
	PublishEvent(subject MessageSubject, content interface{}) error
	OnEvent(subject MessageSubject, f func(event Envelope)) error
}

var impl EventStore

func SetEventStore(es EventStore) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishEvent(subject MessageSubject, content interface{}) error {
	return impl.PublishEvent(subject, content)
}

func OnEvent(subject MessageSubject, f func(event Envelope)) error {
	return impl.OnEvent(subject, f)
}
