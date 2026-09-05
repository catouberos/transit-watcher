package events

import (
	"uuid"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"google.golang.org/protobuf/proto"
)

const (
	ApplicationCloudEventsProtobuf = "application/cloudevents+protobuf"
)

func (c *kafkaEventHandler) CreateEvent(message proto.Message) (cloudevents.Event, error) {
	payload, err := proto.Marshal(message)
	if err != nil {
		return cloudevents.NewEvent(), err
	}

	event := cloudevents.NewEvent()
	event.SetID(uuid.NewV7().String())
	event.SetSource(c.name)
	event.SetSpecVersion(cloudevents.VersionV1)
	event.SetType(string(message.ProtoReflect().Descriptor().Name()))
	event.SetData(ApplicationCloudEventsProtobuf, payload)

	return event, nil
}
