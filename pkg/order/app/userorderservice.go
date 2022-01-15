package app

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/app/integrationevent"
	"store/pkg/common/app/storedevent"
)

var ErrPaymentFailed = errors.New("order payment failed")

type UserOrderService interface {
	Create(userID UserID, price float64) (OrderID, error)
}

func NewUserOrderService(trUnitFactory TransactionalUnitFactory, userOrderReadRepository UserOrderRepository, eventSender storedevent.Sender, billingClient BillingClient) UserOrderService {
	return &userOrderService{
		trUnitFactory:           trUnitFactory,
		userOrderReadRepository: userOrderReadRepository,
		eventSender:             eventSender,
		billingClient:           billingClient,
	}
}

type userOrderService struct {
	trUnitFactory           TransactionalUnitFactory
	userOrderReadRepository UserOrderRepository
	eventSender             storedevent.Sender
	billingClient           BillingClient
}

func (s *userOrderService) Create(userID UserID, price float64) (OrderID, error) {
	id := ID(uuid.NewV1())
	orderID := OrderID(uuid.NewV1())

	order := NewUserOrder(id, userID, orderID, price, Created)
	paymentSucceeded, err := s.billingClient.ProcessOrderPayment(uuid.UUID(userID), price)
	if err != nil {
		return orderID, err
	}

	err = s.executeInTransaction(func(provider RepositoryProvider) error {
		var event integrationevent.EventData
		if paymentSucceeded {
			err2 := provider.UserOrderRepository().Store(order)
			if err2 != nil {
				return err2
			}
			event = NewUserOrderConfirmedEvent(orderID, userID)
		} else {
			event = NewUserOrderRejectedEvent(orderID, userID)
		}

		err2 := provider.EventStore().Add(event)
		if err2 != nil {
			return err2
		}
		s.eventSender.EventStored(event.UID)

		return nil
	})
	if err != nil {
		return orderID, err
	}

	s.eventSender.SendStoredEvents()

	if !paymentSucceeded {
		return orderID, errors.WithStack(ErrPaymentFailed)
	}

	return orderID, err
}

func (s *userOrderService) executeInTransaction(f func(RepositoryProvider) error) (err error) {
	var trUnit TransactionalUnit
	trUnit, err = s.trUnitFactory.NewTransactionalUnit()
	if err != nil {
		return err
	}
	defer func() {
		err = trUnit.Complete(err)
	}()
	err = f(trUnit)
	return err
}
