package messaging

import (
	"context"
	"encoding/json"
	commons "order-service-backend/internal/common"
	"order-service-backend/internal/models"
	"order-service-backend/internal/models/consumer"
	"order-service-backend/internal/usecase"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

type OrderConsumer struct {
	UseCase       *usecase.OrderUseCase
	CreditLimitCh chan []consumer.CreditLimitEvent
}

func NewOrderConsumer(useCase *usecase.OrderUseCase, ch chan []consumer.CreditLimitEvent) *OrderConsumer {
	return &OrderConsumer{
		UseCase:       useCase,
		CreditLimitCh: ch,
	}
}

func (c *OrderConsumer) ConsumeCreditLimitData(log *zap.Logger, events *models.Events) {
	msgs, err := events.SubGroup("check").Subscribe(context.TODO(), commons.CreditLimitExchange)

	if err != nil {
		log.Error(err.Error())
	}

	go func(msgs <-chan *message.Message) {
		for msg := range msgs {
			data := &[]consumer.CreditLimitEvent{}
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				continue
			}

			c.CreditLimitCh <- *data

			msg.Ack()
			log.Info("success consume message topic: " + commons.CreditLimitDataQueue)
		}
	}(msgs)
}
