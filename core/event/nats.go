package event

import (
	"bytes"
	"encoding/gob"
	"github.com/nats-io/go-nats"
	"github.com/plopezm/cloud-kaiser/core/types"
)

type NatsEventStore struct {
	nc                      *nats.Conn
	taskCreatedSubscription *nats.Subscription
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}

func (eventStore *NatsEventStore) Close() {
	if eventStore.nc != nil {
		eventStore.nc.Close()
	}
	if eventStore.taskCreatedSubscription != nil {
		eventStore.taskCreatedSubscription.Unsubscribe()
	}
}

func (eventStore *NatsEventStore) PublishEvent(subject MessageSubject, content interface{}) error {
	return eventStore.publishMessage(Envelope{
		Destination: TaskCreated,
		Content:     content,
	})
}

func (eventStore *NatsEventStore) OnEvent(subject MessageSubject, f func(event Envelope)) error {
	var err error
	content := types.Task{}
	eventStore.taskCreatedSubscription, err = eventStore.nc.Subscribe(string(subject), func(msg *nats.Msg) {
		eventStore.readMessage(msg.Data, &content)
		f(Envelope{
			Destination: MessageSubject(msg.Subject),
			Content:     content,
		})
	})
	return err
}

func (eventStore *NatsEventStore) publishMessage(envelope Envelope) error {
	data, err := eventStore.writeMessage(envelope.Content)
	if err != nil {
		return err
	}
	return eventStore.nc.Publish(string(envelope.Destination), data)
}

func (eventStore *NatsEventStore) writeMessage(content interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(content)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (mq *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}
