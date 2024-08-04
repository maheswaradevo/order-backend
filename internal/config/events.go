package config

import (
	"fmt"
	"order-service-backend/internal/models"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
)

func NewEvent(config *Config) models.Events {
	var address string

	if config.RabbitMqConfig.SSL {
		address = fmt.Sprintf("amqps://%s:%s@%s:%s/", config.RabbitMqConfig.Username, config.RabbitMqConfig.Password, config.RabbitMqConfig.Address, config.RabbitMqConfig.Port)
	} else {
		address = fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitMqConfig.Username, config.RabbitMqConfig.Password, config.RabbitMqConfig.Address, config.RabbitMqConfig.Port)
	}

	amqpConfig := amqp.NewDurableQueueConfig(address)
	brodcastConfig := amqp.NewDurablePubSubConfig(
		address,
		nil,
	)

	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		watermill.NewStdLogger(false, false),
	)

	if err != nil {
		panic(err)
	}

	publisher, err := amqp.NewPublisher(
		amqpConfig,
		watermill.NewStdLogger(false, false),
	)

	if err != nil {
		panic(err)
	}

	brodcastPublisher, err := amqp.NewPublisher(
		brodcastConfig,
		watermill.NewStdLogger(false, false),
	)

	if err != nil {
		panic(err)
	}

	return models.Events{
		Publisher:  publisher,
		Subscriber: subscriber,
		Brodcast:   brodcastPublisher,
		SubGroup: func(prefix string) models.Subscriber {
			subGroup := amqp.NewDurablePubSubConfig(
				address,
				amqp.GenerateQueueNameTopicNameWithSuffix(prefix),
			)

			brodcastSubGroup, err := amqp.NewSubscriber(
				subGroup,
				watermill.NewStdLogger(false, false),
			)

			if err != nil {
				panic(err)
			}

			return brodcastSubGroup
		},
	}
}
