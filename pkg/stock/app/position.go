package app

type Position interface {
	ID() PositionID
	Total() int
}

func NewPosition(id PositionID, total int) Position {
	return &position{
		id:    id,
		total: total,
	}
}

type position struct {
	id    PositionID
	total int
}

func (p *position) ID() PositionID {
	return p.id
}

func (p *position) Total() int {
	return p.total
}
