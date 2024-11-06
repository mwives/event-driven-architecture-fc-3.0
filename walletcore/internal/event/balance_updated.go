package event

import "time"

type BalanceUpdatedEvent struct {
	Name    string
	Payload interface{}
}

func NewBalanceUpdatedEvent() *BalanceUpdatedEvent {
	return &BalanceUpdatedEvent{
		Name: "balance_updated",
	}
}

func (e *BalanceUpdatedEvent) GetName() string {
	return e.Name
}

func (e *BalanceUpdatedEvent) GetDateTime() time.Time {
	return time.Now()
}

func (e *BalanceUpdatedEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *BalanceUpdatedEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}
