package app

type Status int

const (
	Created Status = iota
	Canceled
	Paid
	Processing
	Shipped
	Completed
)
