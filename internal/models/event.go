package models

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type Events struct {
	Publisher  Publisher
	Subscriber Subscriber
	Brodcast   Publisher
	SubGroup   func(prefix string) Subscriber
}

type Publisher interface {
	Close() error
	Publish(topic string, messages ...*message.Message) (err error)
}

type Subscriber interface {
	Close() error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	SubscribeInitialize(topic string) (err error)
}

type Event interface{}
