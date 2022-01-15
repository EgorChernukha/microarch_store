package storedevent

import (
	"context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"store/pkg/common/app/integrationevent"
	"store/pkg/common/app/storedevent"
	"store/pkg/common/app/streams"
	internalintegrationevent "store/pkg/common/infrastructure/integrationevent"
)

const eventDispatchDelay = time.Second * 5
const unconfirmedEventSendDelay = time.Second * 30

func NewEventSender(ctx context.Context, repository storedevent.EventStore, rmqEnv streams.Environment) (storedevent.Sender, error) {
	sender := &sender{
		repository: repository,
		producer:   nil,
	}
	producer, err := rmqEnv.AddProducer(streams.IntegrationEventStreamName, sender.msgConfirmation)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sender.producer = producer

	sender.start(ctx)
	return sender, nil
}

type sender struct {
	repository     storedevent.EventStore
	producer       streams.Producer
	addedEventUids []integrationevent.EventUID
}

func (s *sender) EventStored(uid integrationevent.EventUID) {
	s.addedEventUids = append(s.addedEventUids, uid)
}

func (s *sender) SendStoredEvents() {
	if len(s.addedEventUids) == 0 {
		return
	}
	s.sendEvents(&s.addedEventUids)
	s.addedEventUids = nil
}

func (s *sender) msgConfirmation(id streams.MsgID, confirmed bool, err error) {
	if !confirmed {
		logrus.Errorf("send stream message error - %v", err)
		return
	}
	err = s.repository.ConfirmDelivery(storedevent.EventID(id))
	if err != nil {
		logrus.Error(err)
	}
}

func (s *sender) sendEvents(specificUids *[]integrationevent.EventUID) {
	var events []storedevent.Event
	var err error
	if specificUids != nil {
		events, err = s.repository.FindByUIDs(*specificUids)
	} else {
		events, err = s.repository.FindAllUnconfirmedBefore(time.Now().Add(-unconfirmedEventSendDelay))
	}
	if err != nil {
		logrus.Error(err)
		return
	}
	if len(events) == 0 {
		return
	}

	streamMessages := make([]streams.Msg, 0, len(events))
	for _, event := range events {
		data, err := json.Marshal(internalintegrationevent.EventDataView{
			UID:  uuid.UUID(event.UID).String(),
			Type: event.Type,
			Body: event.Body,
		})
		if err != nil {
			logrus.Error(err)
			continue
		}
		streamMessages = append(streamMessages, streams.Msg{
			ID:   streams.MsgID(event.ID),
			Body: string(data),
		})
	}

	err = s.producer.BatchSend(streamMessages)
	if err != nil {
		logrus.Error(err)
	}
}

func (s *sender) start(ctx context.Context) {
	ticker := time.NewTicker(eventDispatchDelay)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.sendEvents(nil)
			case <-ctx.Done():
				return
			}
		}
	}()
}
