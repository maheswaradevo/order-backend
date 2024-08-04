package messaging

import (
	"encoding/json"

	commons "github.com/maheswaradevo/order-backend/internal/common"
	"github.com/maheswaradevo/order-backend/internal/models"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

type OrderPublisher struct {
	event *models.Events
	log   *zap.Logger
}

func (o *OrderPublisher) PushCreditLimitRequest(data *models.CreditLimitRequest) (bool, error) {
	pyld, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	request := models.MessageRabbit{
		Exchange: commons.CreditLimitExchange,
		Queue:    commons.CreditLimitDataRequest,
		Payload:  pyld,
	}

	if _, err := o.Publish(request); err != nil {
		o.log.Error("failed to publish data: ", zap.Error(err))
		return false, err
	}
	return true, nil
}

func (o *OrderPublisher) PushUpdateCreditLimit(data *models.CreditLimitUpdate) (bool, error) {
	pyld, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	request := models.MessageRabbit{
		Exchange: commons.CreditLimitExchange,
		Queue:    commons.CreditLimitDataUpdate,
		Payload:  pyld,
	}
	if _, err := o.Publish(request); err != nil {
		o.log.Error("failed to publish data: ", zap.Error(err))
		return false, err
	}
	return true, nil
}

func (o *OrderPublisher) Publish(data models.MessageRabbit) (bool, error) {
	if err := o.event.Publisher.Publish(data.Queue, &message.Message{
		Payload: data.Payload,
	}); err != nil {
		o.log.Error("failed to broadcast message to "+data.Queue, zap.Error(err))
		return false, err
	}

	o.log.Info("successfully broadcast message to " + data.Queue)
	return true, nil
}

func NewOrderPublisher(events *models.Events, log *zap.Logger) *OrderPublisher {
	return &OrderPublisher{
		log:   log,
		event: events,
	}
}
