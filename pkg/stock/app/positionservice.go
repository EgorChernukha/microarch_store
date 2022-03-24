package app

import uuid "github.com/satori/go.uuid"

type ReservePositionInputItem struct {
	OrderID    uuid.UUID
	PositionID uuid.UUID
	Count      int
}

type ReservePositionInput struct {
	Items []ReservePositionInputItem
}

type PositionService interface {
	AddPosition(title string, total int) (uuid.UUID, error)
	TopUpPosition(positionID uuid.UUID, count int) error
	ReservePosition(input ReservePositionInput) error
	ConfirmReserves(orderID uuid.UUID) error
	CancelReserves(orderID uuid.UUID) error
}

func NewPositionService(positionRepository PositionRepository, orderPositionRepository OrderPositionRepository) PositionService {
	return &positionService{
		positionRepository:      positionRepository,
		orderPositionRepository: orderPositionRepository,
	}
}

type positionService struct {
	positionRepository      PositionRepository
	orderPositionRepository OrderPositionRepository
}

func (p *positionService) AddPosition(title string, total int) (uuid.UUID, error) {
	id := p.positionRepository.NewID()
	position := NewPosition(id, title, total)

	err := p.positionRepository.Store(position)

	return uuid.UUID(id), err
}

func (p *positionService) TopUpPosition(positionID uuid.UUID, count int) error {
	position, err := p.positionRepository.FindByID(PositionID(positionID))
	if err != nil {
		return err
	}

	err = position.TopUp(count)
	if err != nil {
		return err
	}

	return p.positionRepository.Store(position)
}

func (p *positionService) ReservePosition(input ReservePositionInput) error {
	for _, reservePositionInputItem := range input.Items {
		positionID := PositionID(reservePositionInputItem.PositionID)
		orderPositionID := p.orderPositionRepository.NewID()
		orderPosition := NewOrderPosition(
			orderPositionID,
			OrderID(reservePositionInputItem.OrderID),
			positionID,
			reservePositionInputItem.Count,
			Reserved,
		)

		position, err := p.positionRepository.FindByID(positionID)
		if err != nil {
			return err
		}

		err = position.Sub(reservePositionInputItem.Count)
		if err != nil {
			return err
		}

		err = p.positionRepository.Store(position)
		if err != nil {
			return err
		}

		err = p.orderPositionRepository.Store(orderPosition)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *positionService) ConfirmReserves(orderID uuid.UUID) error {
	orderPositions, err := p.orderPositionRepository.FindByOrderID(OrderID(orderID))
	if err != nil {
		return err
	}

	for _, orderPosition := range orderPositions {
		orderPosition.Confirm()
		err = p.orderPositionRepository.Store(orderPosition)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *positionService) CancelReserves(orderID uuid.UUID) error {
	orderPositions, err := p.orderPositionRepository.FindByOrderID(OrderID(orderID))
	if err != nil {
		return err
	}

	for _, orderPosition := range orderPositions {
		orderPosition.Cancel()

		position, err := p.positionRepository.FindByID(orderPosition.PositionID())
		if err != nil {
			return err
		}

		err = position.TopUp(orderPosition.Count())
		if err != nil {
			return err
		}

		err = p.positionRepository.Store(position)
		if err != nil {
			return err
		}

		err = p.orderPositionRepository.Store(orderPosition)
		if err != nil {
			return err
		}
	}

	return nil
}
