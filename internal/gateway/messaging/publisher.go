package messaging

import (
	"encoding/json"
	"order-service-backend/internal/models"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

type Publisher[T models.Event] struct {
	Publisher models.Publisher
	Queue     string
	Log       *zap.Logger
}

func (p *Publisher[T]) GetQueue() *string {
	return &p.Queue
}

func (p *Publisher[T]) Publish(event T) error {
	value, err := json.Marshal(event)
	if err != nil {
		p.Log.Error("failed to marshal event: ", zap.Error(err))
		return err
	}

	queue := p.GetQueue()
	if err := p.Publisher.Publish(*queue, &message.Message{
		Payload: value,
	}); err != nil {
		p.Log.Error("failed broadcast message to "+*queue, zap.Error(err))
	}

	p.Log.Info("successfully broadcast message to " + *queue)
	return nil
}
