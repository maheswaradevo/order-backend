package messaging

import (
	"context"
	"encoding/json"
	commons "order-service-backend/internal/common"
	"order-service-backend/internal/models"
	"order-service-backend/internal/models/consumer"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func ConsumeUserData(log *zap.Logger, events *models.Events) {
	msgs, err := events.SubGroup("transaction").Subscribe(context.TODO(), commons.UserDataExchange)
	if err != nil {
		log.Error(err.Error())
	}

	go func(msgs <-chan *message.Message) {
		for msg := range msgs {
			data := &consumer.UserEvent{}
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				continue
			}

			msg.Ack()
			log.Info("success consume message topic: " + commons.UserDataSent)
		}
	}(msgs)
}
