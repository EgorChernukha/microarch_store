package app

import uuid "github.com/satori/go.uuid"

type OrderDeliveryService interface {
	AddOrderDelivery(orderID uuid.UUID, userID uuid.UUID) (uuid.UUID, error)
	ConfirmOrderDelivery(orderDeliveryID uuid.UUID) error
	SentOrderDelivery(orderDeliveryID uuid.UUID) error
	ReceiveOrderDelivery(orderDeliveryID uuid.UUID) error
	RejectOrderDelivery(orderDeliveryID uuid.UUID) error
}

func NewOrderDeliveryService(orderDeliveryRepository OrderDeliveryRepository) OrderDeliveryService {
	return &orderDeliveryService{orderDeliveryRepository: orderDeliveryRepository}
}

type orderDeliveryService struct {
	orderDeliveryRepository OrderDeliveryRepository
}

func (o *orderDeliveryService) AddOrderDelivery(orderID uuid.UUID, userID uuid.UUID) (uuid.UUID, error) {
	id := o.orderDeliveryRepository.NewID()

	orderDelivery := NewOrderDelivery(id, OrderID(orderID), UserID(userID), Created)

	err := o.orderDeliveryRepository.Store(orderDelivery)

	return uuid.UUID(id), err
}

func (o *orderDeliveryService) ConfirmOrderDelivery(orderDeliveryID uuid.UUID) error {
	orderDelivery, err := o.orderDeliveryRepository.FindByID(ID(orderDeliveryID))
	if err != nil {
		return err
	}

	orderDelivery.Confirm()

	return nil
}

func (o *orderDeliveryService) SentOrderDelivery(orderDeliveryID uuid.UUID) error {
	orderDelivery, err := o.orderDeliveryRepository.FindByID(ID(orderDeliveryID))
	if err != nil {
		return err
	}

	orderDelivery.Sent()

	return nil
}

func (o *orderDeliveryService) ReceiveOrderDelivery(orderDeliveryID uuid.UUID) error {
	orderDelivery, err := o.orderDeliveryRepository.FindByID(ID(orderDeliveryID))
	if err != nil {
		return err
	}

	orderDelivery.Receive()

	return nil
}

func (o *orderDeliveryService) RejectOrderDelivery(orderDeliveryID uuid.UUID) error {
	orderDelivery, err := o.orderDeliveryRepository.FindByID(ID(orderDeliveryID))
	if err != nil {
		return err
	}

	orderDelivery.Reject()

	return nil
}
