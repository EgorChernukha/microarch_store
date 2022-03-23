package app

import "errors"

var ErrNotEnoughCount = errors.New("not enough total")
var ErrInvalidCount = errors.New("invalid count")

type Position interface {
	ID() PositionID
	Title() string
	Total() int
	TopUp(count int) error
	Sub(count int) error
}

func NewPosition(id PositionID, title string, total int) Position {
	return &position{
		id:    id,
		title: title,
		total: total,
	}
}

type position struct {
	id    PositionID
	title string
	total int
}

func (p *position) ID() PositionID {
	return p.id
}

func (p *position) Title() string {
	return p.title
}

func (p *position) Total() int {
	return p.total
}

func (p *position) TopUp(count int) error {
	if count < 0 {
		return ErrInvalidCount
	}

	p.total += count

	return nil
}

func (p *position) Sub(count int) error {
	if count < 0 {
		return ErrInvalidCount
	}

	if p.total < count {
		return ErrNotEnoughCount
	}

	p.total -= count
	return nil
}
