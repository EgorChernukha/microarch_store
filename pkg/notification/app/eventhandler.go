package app

import (
	"store/pkg/common/app/integrationevent"
	"store/pkg/notification/domain"
)

type IntegrationEventParser interface {
	ParseIntegrationEvent(event integrationevent.EventData) (UserEvent, error)
}

func NewEventHandler(trUnitFactory TransactionalUnitFactory, parser IntegrationEventParser) integrationevent.EventHandler {
	return &eventHandler{
		trUnitFactory: trUnitFactory,
		parser:        parser,
	}
}

type eventHandler struct {
	trUnitFactory TransactionalUnitFactory
	parser        IntegrationEventParser
}

func (handler *eventHandler) Handle(event integrationevent.EventData) error {
	parsedEvent, err := handler.parser.ParseIntegrationEvent(event)
	if err != nil || parsedEvent == nil {
		return err
	}

	return handler.executeInTransaction(func(provider RepositoryProvider) error {
		eventRepo := provider.ProcessedEventRepository()
		alreadyProcessed, err := eventRepo.SetProcessed(event.UID)
		if err != nil {
			return err
		}
		if alreadyProcessed {
			return nil
		}

		switch e := parsedEvent.(type) {
		case orderConfirmedEvent:
			return handleOrderConfirmedEvent(provider, e)
		case orderRejectedEvent:
			return handleOrderRejectedEvent(provider, e)
		default:
			return nil
		}
	})
}

func (handler *eventHandler) executeInTransaction(f func(RepositoryProvider) error) (err error) {
	var trUnit TransactionalUnit
	trUnit, err = handler.trUnitFactory.NewTransactionalUnit()
	if err != nil {
		return err
	}
	defer func() {
		err = trUnit.Complete(err)
	}()
	err = f(trUnit)
	return err
}

func handleOrderConfirmedEvent(provider RepositoryProvider, e orderConfirmedEvent) error {
	userNotificationDomainService := domain.NewUserNotificationService(provider.UserNotificationRepository())
	userNotificationAppService := NewUserNotificationService(userNotificationDomainService)
	return userNotificationAppService.CreateUserNotification(TypeOrderConfirmed, e.UserID(), e.orderID)
}

func handleOrderRejectedEvent(provider RepositoryProvider, e orderRejectedEvent) error {
	userNotificationDomainService := domain.NewUserNotificationService(provider.UserNotificationRepository())
	userNotificationAppService := NewUserNotificationService(userNotificationDomainService)
	return userNotificationAppService.CreateUserNotification(TypeOrderRejected, e.UserID(), e.orderID)
}
