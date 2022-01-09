package integrationevent

import (
	"github.com/sirupsen/logrus"

	"store/pkg/order/app/integrationevent"
)

func NewIntegrationEventHandler(eventHandlers []integrationevent.Handler) Handler {
	return &handler{eventHandlers: eventHandlers}
}

type handler struct {
	eventHandlers []integrationevent.Handler
}

func (h *handler) Handle(msg string) error {
	event, err := h.parseClientMessage(msg)
	if err != nil {
		logrus.Error(err, "failed to parse integration event")
		return err
	}

	// Skip unsupported events. If you need handle event just add parsing for it.
	if event == nil {
		return nil
	}

	for _, integrationEventHandler := range h.eventHandlers {
		err = integrationEventHandler.Handle(event)
		if err != nil {
			break
		}
	}

	if err != nil {
		logrus.WithField("event_body", msg).Error(err, "failed to process integration event")
	} else {
		logrus.WithField("event_body", msg).Info("integration event received")
	}
	return err
}
