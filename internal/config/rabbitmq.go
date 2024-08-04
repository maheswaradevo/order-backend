package config

import (
	"fmt"
	"time"

	commons "github.com/maheswaradevo/order-backend/internal/common"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	URL        string
	Conn       *amqp.Connection
	Chann      *amqp.Channel
	CloseChann chan *amqp.Error
	QuitChann  chan bool
}

func NewRabbitMQ(option RabbitMQConfig) (*RabbitMQ, error) {
	var addr string

	if option.SSL {
		addr = fmt.Sprintf("amqps://%s:%s@%s:%s/", option.Username, option.Password, option.Address, option.Port)
	} else {
		addr = fmt.Sprintf("amqp://%s:%s@%s:%s/", option.Username, option.Password, option.Address, option.Port)
	}

	rmq := &RabbitMQ{
		URL: addr,
	}

	err := rmq.load()
	if err != nil {
		return nil, err
	}

	rmq.QuitChann = make(chan bool)

	go rmq.handleDisconnect()

	return rmq, nil
}

func (rmq *RabbitMQ) load() error {
	var err error

	rmq.Conn, err = amqp.Dial(rmq.URL)
	if err != nil {
		return err
	}

	rmq.Chann, err = rmq.Conn.Channel()
	if err != nil {
		return err
	}

	rmq.CloseChann = make(chan *amqp.Error)
	rmq.Conn.NotifyClose(rmq.CloseChann)

	// make exchage if not exist
	for _, queue := range commons.QueueArr {
		err := rmq.makeExchange(queue.Exchange, queue.UseDelay)
		if err != nil {
			return err
		}

		err = rmq.makeQueue(queue.Name, queue.RoutingKey, queue.Exchange)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rmq *RabbitMQ) makeExchange(name string, isDelay bool) error {
	var err error

	if isDelay {
		args := make(amqp.Table)
		args["x-delayed-type"] = "direct"
		err = rmq.Chann.ExchangeDeclare(name, "x-delayed-message", true, false, false, false, args)
		if err != nil {
			return fmt.Errorf("%v declaring exchange with delay %v", err, name)
		}
	}

	if !isDelay {
		err = rmq.Chann.ExchangeDeclare(name, "fanout", true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("%v declaring exchange %v", err, name)
		}
	}

	return err
}

func (rmq *RabbitMQ) makeQueue(name string, key string, exchange string) error {
	queue, err := rmq.Chann.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("%v declaring queue %v", err, name)
	}

	err = rmq.Chann.QueueBind(queue.Name, key, exchange, false, nil)
	if err != nil {
		return fmt.Errorf("%v binding queue %v", err, name)
	}

	return err
}

// handleDisconnect will handle disconnect from server and try every 5 second
func (rmq *RabbitMQ) handleDisconnect() {
	for {
		select {
		case errChann := <-rmq.CloseChann:
			if errChann != nil {
				zap.L().Error("rabbitMQ disconnection: %v", zap.Error(errChann))
			}
		case <-rmq.QuitChann:
			rmq.Conn.Close()
			rmq.QuitChann <- true
			return
		}

		zap.L().Info("...trying to reconnect to rabbitMQ...")

		time.Sleep(5 * time.Second)

		if err := rmq.load(); err != nil {
			zap.L().Error("rabbitMQ error: %v", zap.Error(err))
		}
	}
}
