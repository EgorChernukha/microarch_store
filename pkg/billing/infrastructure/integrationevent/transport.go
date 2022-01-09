package integrationevent

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	commonamqp "store/pkg/common/infrastructure/amqp"
)

const (
	domainEventsQueueName    = "billing_domain_event"
	domainEventsExchangeName = "domain_event"
	domainEventsExchangeType = "topic"
	routingKey               = "#"
	routingPrefix            = "billing."
	contentType              = "application/json; charset=utf-8"
)

type Handler interface {
	Handle(msg string) error
}

type Transport interface {
	commonamqp.Channel
	SetHandler(handler Handler)
}

type transport struct {
	conn                  *amqp.Connection
	writeChannel          *amqp.Channel
	handler               Handler
	suppressEventsReading bool
}

func (t *transport) Name() string {
	return "amqp_integration_events"
}

func (t *transport) Send(msgBody, eventType string) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  contentType,
		Body:         []byte(msgBody),
	}
	routingKey := routingPrefix + eventType
	return t.writeChannel.Publish(domainEventsExchangeName, routingKey, false, false, msg)
}

func (t *transport) Connect(conn *amqp.Connection) error {
	t.writeChannel = nil

	t.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	t.writeChannel = channel

	err = channel.ExchangeDeclare(domainEventsExchangeName, domainEventsExchangeType, true, false, false, false, nil)
	if err != nil {
		return err
	}

	if !t.suppressEventsReading {
		return t.connectReadChannel(err, channel)
	}
	return nil
}

func (t *transport) connectReadChannel(err error, channel *amqp.Channel) error {
	readQueue, err := channel.QueueDeclare(domainEventsQueueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = channel.QueueBind(readQueue.Name, routingKey, domainEventsExchangeName, false, nil)
	if err != nil {
		return err
	}

	if t.handler == nil {
		return errors.New("event handler should be set before consuming messages")
	}

	readChan, err := channel.Consume(readQueue.Name, "", false, false, false, false, nil)

	go func() {
		for msg := range readChan {
			err = t.handler.Handle(string(msg.Body))
			if err == nil {
				err = msg.Ack(false)
			} else {
				err = msg.Nack(false, true)
			}
			_ = err
		}
	}()

	return err
}

func (t *transport) SetHandler(handler Handler) {
	t.handler = handler
}

func NewIntegrationEventsTransport(suppressEventsReading bool) Transport {
	return &transport{suppressEventsReading: suppressEventsReading}
}
