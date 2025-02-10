package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/iRankHub/backend/internal/services/notification/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQStorage struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	exchange      string
	queueName     string
	routingKey    string
	metadataStore MetadataStore
	mu            sync.RWMutex
	isConnected   bool
}

func NewRabbitMQStorage(url string, metadataStore MetadataStore) (*RabbitMQStorage, error) {
	s := &RabbitMQStorage{
		exchange:      "notifications",
		queueName:     "notifications.store",
		routingKey:    "notification.#",
		metadataStore: metadataStore,
	}

	if err := s.connect(url); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *RabbitMQStorage) connect(url string) error {
	conn, err := amqp.Dial(url)
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
		s.exchange,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue
	_, err = ch.QueueDeclare(
		s.queueName,
		true,  // durable
		false, // auto-deleted
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-message-ttl": 30 * 24 * 60 * 60 * 1000, // 30 days in milliseconds
		},
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		s.queueName,
		s.routingKey,
		s.exchange,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	s.conn = conn
	s.channel = ch
	s.isConnected = true

	return nil
}

func (s *RabbitMQStorage) Store(ctx context.Context, notification *models.Notification) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	// Store metadata first
	if err := s.metadataStore.StoreMetadata(ctx, notification); err != nil {
		return fmt.Errorf("failed to store metadata: %w", err)
	}

	// Marshal notification
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Create routing key
	routingKey := fmt.Sprintf("notification.%s.%s", notification.Category, notification.Type)

	// Publish to RabbitMQ
	err = s.channel.PublishWithContext(
		ctx,
		s.exchange,
		routingKey,
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Headers: amqp.Table{
				"user_id":    notification.UserID,
				"category":   string(notification.Category),
				"type":       string(notification.Type),
				"expires_at": notification.ExpiresAt.Unix(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (s *RabbitMQStorage) Get(ctx context.Context, notificationID string) (*models.Notification, error) {
	// Get message from RabbitMQ
	msg, ok, err := s.channel.Get(
		s.queueName,
		false, // auto-ack
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("notification not found")
	}

	var notification models.Notification
	if err := json.Unmarshal(msg.Body, &notification); err != nil {
		msg.Nack(false, true) // requeue message
		return nil, fmt.Errorf("failed to unmarshal notification: %w", err)
	}

	msg.Ack(false)
	return &notification, nil
}

func (s *RabbitMQStorage) GetUnread(ctx context.Context, userID string) ([]*models.Notification, error) {
	// This implementation will need to be enhanced based on your specific requirements
	// Currently, it's just getting messages from the queue without filtering
	var notifications []*models.Notification

	msgs, err := s.channel.Consume(
		s.queueName,
		"",    // consumer
		false, // auto-ack
		true,  // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	// Add timeout to prevent blocking indefinitely
	timeout := time.After(5 * time.Second)
	msgChan := make(chan *models.Notification)
	errChan := make(chan error)

	go func() {
		for msg := range msgs {
			var notification models.Notification
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				msg.Nack(false, true) // requeue message
				errChan <- fmt.Errorf("failed to unmarshal notification: %w", err)
				return
			}

			if notification.UserID == userID && !notification.IsRead {
				msg.Ack(false)
				msgChan <- &notification
			} else {
				msg.Nack(false, true) // requeue message
			}
		}
	}()

	select {
	case <-timeout:
		return notifications, nil
	case err := <-errChan:
		return nil, err
	case notification := <-msgChan:
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (s *RabbitMQStorage) MarkAsRead(ctx context.Context, notificationID string) error {
	// Update metadata in database
	metadata := map[string]interface{}{
		"is_read": true,
		"read_at": time.Now(),
	}
	return s.metadataStore.UpdateMetadata(ctx, notificationID, metadata)
}

func (s *RabbitMQStorage) UpdateDeliveryStatus(ctx context.Context, notificationID string, method models.DeliveryMethod, status models.Status) error {
	metadata := map[string]interface{}{
		"delivery_status": map[string]interface{}{
			string(method): map[string]interface{}{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	}
	return s.metadataStore.UpdateMetadata(ctx, notificationID, metadata)
}

func (s *RabbitMQStorage) DeleteExpired(ctx context.Context) error {
	// RabbitMQ handles message expiration automatically via TTL
	// We just need to clean up the metadata
	return s.metadataStore.DeleteExpiredMetadata(ctx)
}

func (s *RabbitMQStorage) GetForRetry(ctx context.Context) ([]*models.Notification, error) {
	// This would typically involve checking the metadata store for failed notifications
	// and then retrieving them from RabbitMQ
	// Implementation depends on your specific retry requirements
	return nil, nil
}

func (s *RabbitMQStorage) Close() error {
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
