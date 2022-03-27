package app

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/app/integrationevent"
	"store/pkg/common/app/storedevent"
)

var ErrPaymentFailed = errors.New("order payment failed")
var ErrReserveOrderDeliveryFailed = errors.New("reserve order delivery failed")
var ErrReserveOrderPositionsFailed = errors.New("reserve order positions failed")

type UserOrderService interface {
	Create(userID UserID, price float64, positionID PositionID, count int) (OrderID, error)
}

func NewUserOrderService(
	trUnitFactory TransactionalUnitFactory,
	userOrderReadRepository UserOrderRepository,
	eventSender storedevent.Sender,
	billingClient BillingClient,
	deliveryClient DeliveryClient,
	stockClient StockClient,
) UserOrderService {
	return &userOrderService{
		trUnitFactory:           trUnitFactory,
		userOrderReadRepository: userOrderReadRepository,
		eventSender:             eventSender,
		billingClient:           billingClient,
		deliveryClient:          deliveryClient,
		stockClient:             stockClient,
	}
}

type userOrderService struct {
	trUnitFactory           TransactionalUnitFactory
	userOrderReadRepository UserOrderRepository
	eventSender             storedevent.Sender
	billingClient           BillingClient
	deliveryClient          DeliveryClient
	stockClient             StockClient
}

func (s *userOrderService) Create(userID UserID, price float64, positionID PositionID, count int) (OrderID, error) {
	id := ID(uuid.NewV1())
	orderID := OrderID(uuid.NewV1())

	order := NewUserOrder(id, userID, orderID, price, Created)
	paymentSucceeded, err := s.billingClient.ProcessOrderPayment(uuid.UUID(userID), price)
	if err != nil {
		return orderID, err
	}

	reserveDeliverySucceeded, err := s.deliveryClient.ReserveDelivery(uuid.UUID(userID), uuid.UUID(orderID))
	if err != nil {
		return orderID, err
	}

	positions := []ReserveOrderPositionInputItem{
		{
			PositionID: uuid.UUID(positionID),
			OrderID:    uuid.UUID(orderID),
			Count:      count,
		},
	}
	reservePositionsSucceeded, err := s.stockClient.ReserveOrderPositions(ReserveOrderPositionInput{Positions: positions})
	if err != nil {
		return orderID, err
	}

	err = s.executeInTransaction(func(provider RepositoryProvider) error {
		var event integrationevent.EventData
		if paymentSucceeded && reserveDeliverySucceeded && reservePositionsSucceeded {
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

	if !reserveDeliverySucceeded {
		return orderID, errors.WithStack(ErrReserveOrderDeliveryFailed)
	}

	if !reservePositionsSucceeded {
		return orderID, errors.WithStack(ErrReserveOrderPositionsFailed)
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
