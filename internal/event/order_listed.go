package event

import "time"

type OrderListed struct {
	Name    string
	Payload interface{}
}

func NewOrderListed() *OrderListed {
	return &OrderListed{
		Name: "OrderListed",
	}
}

func (e *OrderListed) GetName() string {
	return e.Name
}

func (e *OrderListed) GetPayload() interface{} {
	return e.Payload
}

func (e *OrderListed) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *OrderListed) GetDateTime() time.Time {
	return time.Now()
}
