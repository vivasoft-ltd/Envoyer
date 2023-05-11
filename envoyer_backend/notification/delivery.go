package notification

import "time"

type Acknowledger interface {
	Ack(tag uint64, multiple bool) error
	Reject(tag uint64, requeue bool) error
}

type Delivery struct {
	Acknowledger Acknowledger
	Id           uint64
	MessageType  string
	Exchange     string
	RoutingKey   string
	Timestamp    time.Time
	Headers      map[string]interface{}
	Body         []byte
}

type ValidityChecker interface {
	IsValid(delivery Delivery) bool
}
