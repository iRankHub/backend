package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	done          chan bool
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
	isConnected   bool
	url           string
	exchange      string
	queue         string
	routingKey    string
	mutex         sync.Mutex
}

func NewRabbitMQ(url, exchange, queue, routingKey string) (*RabbitMQ, error) {
	rabbitmq := &RabbitMQ{
		done:       make(chan bool),
		url:        url,
		exchange:   exchange,
		queue:      queue,
		routingKey: routingKey,
	}

	go rabbitmq.handleReconnect()

	return rabbitmq, nil
}

func (r *RabbitMQ) handleReconnect() {
	for {
		r.isConnected = false
		log.Println("Attempting to connect to RabbitMQ...")

		conn, err := r.connect()
		if err != nil {
			log.Println("Failed to connect. Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		r.changeConnection(conn)
		log.Println("Connected to RabbitMQ")

		select {
		case <-r.done:
			return
		case <-r.notifyClose:
			log.Println("Connection closed. Reconnecting...")
		}
	}
}

func (r *RabbitMQ) connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		r.exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		r.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		r.queue,
		r.routingKey,
		r.exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (r *RabbitMQ) changeConnection(connection *amqp.Connection) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.conn = connection
	r.notifyClose = make(chan *amqp.Error)
	r.conn.NotifyClose(r.notifyClose)

	r.channel, _ = r.conn.Channel()
	r.notifyConfirm = make(chan amqp.Confirmation, 1)
	r.channel.NotifyPublish(r.notifyConfirm)

	r.isConnected = true
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, body []byte) error {
	if !r.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	err := r.channel.PublishWithContext(
		ctx,
		r.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	select {
	case confirm := <-r.notifyConfirm:
		if confirm.Ack {
			return nil
		}
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
	}

	return fmt.Errorf("message not confirmed")
}

func (r *RabbitMQ) StartConsuming(ctx context.Context, handler func(Notification) error, concurrency int) error {
	for {
		if !r.isConnected {
			log.Println("Waiting for RabbitMQ connection...")
			time.Sleep(1 * time.Second)
			continue
		}

		msgs, err := r.channel.Consume(
			r.queue,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("Failed to consume: %v", err)
			continue
		}

		var wg sync.WaitGroup
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case msg, ok := <-msgs:
						if !ok {
							return
						}
						var notification Notification
						err := json.Unmarshal(msg.Body, &notification)
						if err != nil {
							log.Printf("Error unmarshalling notification: %v", err)
							msg.Nack(false, false)
							continue
						}

						err = handler(notification)
						if err != nil {
							log.Printf("Error handling notification: %v", err)
							msg.Nack(false, true)
						} else {
							msg.Ack(false)
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		}

		wg.Wait()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

func (r *RabbitMQ) Close() error {
	if !r.isConnected {
		return nil
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	close(r.done)

	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}

	r.isConnected = false
	return nil
}
