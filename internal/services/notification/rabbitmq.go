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
	reconnectChan chan struct{} // Channel to signal successful reconnection
}

func NewRabbitMQ(url, exchange, queue, routingKey string) (*RabbitMQ, error) {
	rabbitmq := &RabbitMQ{
		done:          make(chan bool),
		reconnectChan: make(chan struct{}, 1),
		url:           url,
		exchange:      exchange,
		queue:         queue,
		routingKey:    routingKey,
	}

	// Start the reconnection handler
	go rabbitmq.handleReconnect()

	// Wait for the initial connection before returning
	select {
	case <-rabbitmq.reconnectChan:
		log.Println("Successfully connected to RabbitMQ")
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("failed to connect to RabbitMQ within timeout")
	}

	return rabbitmq, nil
}

func (r *RabbitMQ) handleReconnect() {
	for {
		r.setConnected(false)
		log.Println("Attempting to connect to RabbitMQ...")

		conn, err := amqp.Dial(r.url)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Failed to open a channel: %v. Retrying in 5 seconds...", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// Setup the exchange
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
			log.Printf("Failed to declare exchange: %v. Retrying in 5 seconds...", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// Setup the queue
		_, err = ch.QueueDeclare(
			r.queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("Failed to declare queue: %v. Retrying in 5 seconds...", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// Bind the queue to the exchange
		err = ch.QueueBind(
			r.queue,
			r.routingKey,
			r.exchange,
			false,
			nil,
		)
		if err != nil {
			log.Printf("Failed to bind queue: %v. Retrying in 5 seconds...", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// Enable publishing confirmations
		err = ch.Confirm(false)
		if err != nil {
			log.Printf("Failed to enable publish confirmations: %v. Retrying in 5 seconds...", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		r.changeConnection(conn, ch)
		log.Println("Connected to RabbitMQ successfully")

		// Signal that we're connected
		select {
		case r.reconnectChan <- struct{}{}:
		default:
		}

		// Wait for connection to close
		select {
		case <-r.done:
			return
		case err := <-r.notifyClose:
			log.Printf("Connection closed: %v. Reconnecting...", err)
		}
	}
}

func (r *RabbitMQ) setConnected(connected bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.isConnected = connected
}

func (r *RabbitMQ) changeConnection(conn *amqp.Connection, ch *amqp.Channel) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Close existing channel and connection if they exist
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}

	r.conn = conn
	r.channel = ch

	// Set up notification channels
	r.notifyClose = make(chan *amqp.Error)
	r.conn.NotifyClose(r.notifyClose)

	r.notifyConfirm = make(chan amqp.Confirmation, 1)
	r.channel.NotifyPublish(r.notifyConfirm)

	r.isConnected = true
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, body []byte) error {
	// Check connection
	if !r.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Ensure channel is still open
	if r.channel == nil {
		return fmt.Errorf("channel is not open")
	}

	// Publish message
	err := r.channel.PublishWithContext(
		ctx,
		r.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // Make messages persistent
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// Wait for confirmation
	select {
	case confirm := <-r.notifyConfirm:
		if confirm.Ack {
			return nil
		}
		return fmt.Errorf("message not acknowledged")
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return fmt.Errorf("confirmation timeout")
	}
}

func (r *RabbitMQ) StartConsuming(ctx context.Context, handler func(Notification) error, concurrency int) error {
	for {
		// Ensure we're connected
		if !r.isConnected {
			log.Println("Not connected to RabbitMQ. Waiting for connection...")
			select {
			case <-r.reconnectChan:
				log.Println("Connection reestablished, resuming consuming")
			case <-time.After(30 * time.Second):
				log.Println("Failed to start consuming messages: Exception (504) Reason: \"channel/connection is not open\". Retrying in 30s...")
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// Create a local copy of the channel to use for consuming
		r.mutex.Lock()
		if r.channel == nil {
			r.mutex.Unlock()
			log.Println("Channel is nil. Waiting for connection...")
			time.Sleep(5 * time.Second)
			continue
		}
		channel := r.channel
		r.mutex.Unlock()

		// Start consuming
		msgs, err := channel.Consume(
			r.queue,
			"",    // Consumer tag
			false, // Auto-ack
			false, // Exclusive
			false, // No local
			false, // No wait
			nil,   // Args
		)
		if err != nil {
			log.Printf("Failed to start consuming messages: %v. Retrying in 30s...", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Process messages
		var wg sync.WaitGroup
		consumeErrChan := make(chan error, 1)
		stopWorkers := make(chan struct{})

		// Start worker goroutines
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				log.Printf("Starting consumer worker %d", workerID)

				for {
					select {
					case msg, ok := <-msgs:
						if !ok {
							log.Printf("Worker %d: message channel closed", workerID)
							return
						}

						var notification Notification
						err := json.Unmarshal(msg.Body, &notification)
						if err != nil {
							log.Printf("Worker %d: Error unmarshalling notification: %v", workerID, err)
							msg.Nack(false, false) // Don't requeue malformed messages
							continue
						}

						err = handler(notification)
						if err != nil {
							log.Printf("Worker %d: Error handling notification: %v", workerID, err)
							msg.Nack(false, true) // Requeue on handler error
						} else {
							msg.Ack(false)
						}

					case <-stopWorkers:
						log.Printf("Worker %d: stopping", workerID)
						return

					case <-ctx.Done():
						log.Printf("Worker %d: context canceled", workerID)
						return
					}
				}
			}(i)
		}

		// Wait for either context cancellation or connection error
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping consumers")
			close(stopWorkers)
			wg.Wait()
			return ctx.Err()

		case <-r.notifyClose:
			log.Println("RabbitMQ connection closed, restarting consumers")
			close(stopWorkers)
			wg.Wait()
			// Connection will be reestablished by handleReconnect
			continue

		case err := <-consumeErrChan:
			log.Printf("Consumer error: %v, restarting consumers", err)
			close(stopWorkers)
			wg.Wait()
			continue
		}
	}
}

func (r *RabbitMQ) Close() error {
	log.Println("Closing RabbitMQ connection...")

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Signal the reconnect goroutine to stop
	close(r.done)

	// Close channel and connection
	var err error
	if r.channel != nil {
		err = r.channel.Close()
		r.channel = nil
	}

	if r.conn != nil {
		connErr := r.conn.Close()
		if err == nil {
			err = connErr
		}
		r.conn = nil
	}

	r.isConnected = false
	return err
}
