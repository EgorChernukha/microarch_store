package integrationevent

import (
	"encoding/json"

	"github.com/cenkalti/backoff"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"store/pkg/common/app/integrationevent"
	"store/pkg/common/app/streams"
)

func StartEventConsumer(rmqEnv streams.Environment, handler integrationevent.EventHandler) error {
	consumer, err := rmqEnv.AddConsumer(streams.IntegrationEventStreamName)
	if err != nil {
		return err
	}
	eventConsumer := &eventConsumer{
		consumer: consumer,
		handler:  handler,
	}
	consumer.SetMessageHandler(eventConsumer.messageHandler)

	return nil
}

type eventConsumer struct {
	consumer streams.Consumer
	handler  integrationevent.EventHandler
}

func (ec *eventConsumer) messageHandler(msg *amqp.Message) {
	data := msg.Data
	if len(data) == 0 {
		logrus.Error("received message without body")
		return
	}
	if len(data) > 1 {
		logrus.Warnf("received data with multiple data - %v", data)
	}
	rawData := data[0]

	var eventData EventDataView
	err := json.Unmarshal(rawData, &eventData)
	if err != nil {
		logrus.Error("unsupported message body")
		return
	}
	eventID, err := uuid.FromString(eventData.UID)
	if err != nil {
		logrus.Error("invalid event uid")
		return
	}

	ec.handleEvent(integrationevent.EventData{
		UID:  integrationevent.EventUID(eventID),
		Type: eventData.Type,
		Body: eventData.Body,
	})
}

func (ec *eventConsumer) handleEvent(eventData integrationevent.EventData) {
	err := backoff.Retry(func() error {
		err2 := ec.handler.Handle(eventData)
		if err2 != nil {
			logrus.Errorf("error processing integration event - '%s'. attempt to retry", err2.Error())
		}
		return err2
	}, backoff.NewExponentialBackOff())

	if err != nil {
		logrus.Fatalf("error processing integration event - %s\nDetails: uid - '%s', type - '%s', body - '%s'", err.Error(), eventData.UID, eventData.Type, eventData.Body)
	} else {
		logrus.Infof("integration event '%s' with type '%s' handled", uuid.UUID(eventData.UID).String(), eventData.Type)
	}
}
