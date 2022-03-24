package app

import uuid "github.com/satori/go.uuid"

type OrderDeliveryService interface {
	AddOrderDelivery(orderID uuid.UUID, userID uuid.UUID) (uuid.UUID, error)
	MarkOrderDeliveryAsSent(orderDeliveryID uuid.UUID) error
	MarkOrderDeliveryAsReceived(orderDeliveryID uuid.UUID) error
	MarkOrderDeliveryAsRejected(orderDeliveryID uuid.UUID) error
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

func (o *orderDeliveryService) MarkOrderDeliveryAsSent(orderDeliveryID uuid.UUID) error {
	orderDelivery, err := o.orderDeliveryRepository.FindByID(ID(orderDeliveryID))
	if err != nil {
		return err
	}

	orderDelivery.MarkAsSent()

	return nil
}

func (o *orderDeliveryService) MarkOrderDeliveryAsReceived(orderDeliveryID uuid.UUID) error {
	orderDelivery, err := o.orderDeliveryRepository.FindByID(ID(orderDeliveryID))
	if err != nil {
		return err
	}

	orderDelivery.MarkAsReceived()

	return nil
}

func (o *orderDeliveryService) MarkOrderDeliveryAsRejected(orderDeliveryID uuid.UUID) error {
	orderDelivery, err := o.orderDeliveryRepository.FindByID(ID(orderDeliveryID))
	if err != nil {
		return err
	}

	orderDelivery.MarkAsRejected()

	return nil
}
