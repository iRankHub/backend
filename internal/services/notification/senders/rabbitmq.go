package senders

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/iRankHub/backend/internal/services/notification/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQSender struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	config        RabbitMQConfig
	isConnected   bool
	mu            sync.RWMutex
	closeC        chan struct{}
	reconnectC    chan struct{}
	subscriptions map[string]subscriptionInfo
	subMu         sync.RWMutex
}

type subscriptionInfo struct {
	opts       models.SubscriptionOptions
	notifChan  chan *models.Notification
	cancelFunc func()
}

var _ models.Queue = (*RabbitMQSender)(nil)

func NewRabbitMQSender(config RabbitMQConfig) (*RabbitMQSender, error) {
	if config.TTL == 0 {
		config.TTL = 30 * 24 * 60 * 60 * 1000 // 30 days in milliseconds
	}

	sender := &RabbitMQSender{
		config:        config,
		closeC:        make(chan struct{}),
		reconnectC:    make(chan struct{}, 1),
		subscriptions: make(map[string]subscriptionInfo),
	}

	if err := sender.connect(); err != nil {
		return nil, err
	}

	go sender.reconnectLoop()

	return sender, nil
}

func (s *RabbitMQSender) connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conn, err := amqp.Dial(s.config.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		s.config.Exchange,
		s.config.ExchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		amqp.Table{
			"x-message-ttl": s.config.TTL,
		},
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	s.conn = conn
	s.channel = ch
	s.isConnected = true

	// Set up connection monitoring
	go func() {
		<-conn.NotifyClose(make(chan *amqp.Error))
		s.mu.Lock()
		s.isConnected = false
		s.mu.Unlock()
		select {
		case s.reconnectC <- struct{}{}:
		default:
		}
	}()

	return nil
}

func (s *RabbitMQSender) reconnectLoop() {
	for {
		select {
		case <-s.closeC:
			return
		case <-s.reconnectC:
			for !s.isConnected {
				if err := s.connect(); err != nil {
					log.Printf("Failed to reconnect to RabbitMQ: %v", err)
					time.Sleep(5 * time.Second)
					continue
				}
				s.resubscribeAll()
			}
		}
	}
}

func (s *RabbitMQSender) resubscribeAll() {
	s.subMu.RLock()
	defer s.subMu.RUnlock()

	for queueName, subInfo := range s.subscriptions {
		// Create new notification channel
		newChan := make(chan *models.Notification, subInfo.opts.BatchSize)

		// Try to resubscribe
		if _, cleanup, err := s.Subscribe(context.Background(), subInfo.opts); err != nil {
			log.Printf("Failed to resubscribe to queue %s: %v", queueName, err)
			continue
		} else {
			// Update subscription info with new channel and cleanup function
			s.subMu.Lock()
			if oldInfo, exists := s.subscriptions[queueName]; exists {
				// Close old channel and cleanup
				close(oldInfo.notifChan)
				if oldInfo.cancelFunc != nil {
					oldInfo.cancelFunc()
				}

				// Update with new info
				s.subscriptions[queueName] = subscriptionInfo{
					opts:       subInfo.opts,
					notifChan:  newChan,
					cancelFunc: cleanup,
				}
			}
			s.subMu.Unlock()
		}
	}
}

func (s *RabbitMQSender) Publish(ctx context.Context, notification *models.Notification) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	routingKey := fmt.Sprintf("%s.%s.%s",
		notification.UserRole,
		notification.Category,
		notification.Type)

	err = s.channel.PublishWithContext(
		ctx,
		s.config.Exchange,
		routingKey,
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Expiration:   fmt.Sprintf("%d", s.config.TTL),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (s *RabbitMQSender) Subscribe(ctx context.Context, opts models.SubscriptionOptions) (<-chan *models.Notification, func(), error) {
	s.mu.RLock()
	if !s.isConnected {
		s.mu.RUnlock()
		return nil, nil, fmt.Errorf("not connected to RabbitMQ")
	}
	s.mu.RUnlock() // Unlock after checking

	notifChan := make(chan *models.Notification, opts.BatchSize)
	queueName := fmt.Sprintf("notifications.%s.%s", opts.UserRole, opts.UserID)

	// Declare queue
	q, err := s.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-deleted
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-message-ttl": s.config.TTL,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Create binding for each notification type
	for _, t := range opts.Types {
		routingKey := fmt.Sprintf("%s.%s.%s", opts.UserRole, opts.Category, t)
		err = s.channel.QueueBind(
			q.Name,
			routingKey,
			s.config.Exchange,
			false,
			nil,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	msgs, err := s.channel.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack
		true,  // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register consumer: %w", err)
	}

	// Start goroutine to handle messages
	go func() {
		defer close(notifChan)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				var notification models.Notification
				if err := json.Unmarshal(msg.Body, &notification); err != nil {
					log.Printf("Failed to unmarshal notification: %v", err)
					msg.Nack(false, false)
					continue
				}

				select {
				case notifChan <- &notification:
					msg.Ack(false)
				default:
					// Channel is full, nack message and requeue
					msg.Nack(false, true)
				}
			}
		}
	}()

	cleanup := func() {
		s.subMu.Lock()
		defer s.subMu.Unlock()

		if sub, exists := s.subscriptions[queueName]; exists {
			close(sub.notifChan)
			delete(s.subscriptions, queueName)
		}

		if s.channel != nil {
			_ = s.channel.Cancel("", false)
			// Remove queue bindings
			for _, t := range opts.Types {
				routingKey := fmt.Sprintf("%s.%s.%s", opts.UserRole, opts.Category, t)
				_ = s.channel.QueueUnbind(
					queueName,
					routingKey,
					s.config.Exchange,
					nil,
				)
			}
		}
	}

	// Store subscription info
	s.subMu.Lock()
	s.subscriptions[queueName] = subscriptionInfo{
		opts:       opts,
		notifChan:  notifChan,
		cancelFunc: cleanup,
	}
	s.subMu.Unlock()

	return notifChan, cleanup, nil
}

func (s *RabbitMQSender) Unsubscribe(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.channel != nil {
		return s.channel.Cancel("", false)
	}
	return nil
}

func (s *RabbitMQSender) DeleteExpired(ctx context.Context) error {
	// RabbitMQ handles message expiration automatically via TTL
	return nil
}

func (s *RabbitMQSender) Close() error {
	close(s.closeC)

	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	if s.channel != nil {
		err = s.channel.Close()
	}
	if s.conn != nil {
		if cErr := s.conn.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}

	s.isConnected = false
	return err
}
