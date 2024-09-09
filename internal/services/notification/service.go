package notification

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"

	"github.com/iRankHub/backend/internal/models"
)

type NotificationType string

const (
	EmailNotification NotificationType = "email"
	InAppNotification NotificationType = "in_app"
)

type Notification struct {
	Type    NotificationType `json:"type"`
	To      string           `json:"to"`
	Subject string           `json:"subject"`
	Content string           `json:"content"`
}

type NotificationService struct {
	db           *sql.DB
	queries      *models.Queries
	conn         *amqp.Connection
	channel      *amqp.Channel
	queue        amqp.Queue
	emailSender  EmailSender
	inAppSender  InAppSender
}

func NewNotificationService(db *sql.DB) (*NotificationService, error) {
	queries := models.New(db)
	url := viper.GetString("RABBITMQ_URL")
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		"notifications", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	return &NotificationService{
		db:           db,
		queries:      queries,
		conn:         conn,
		channel:      ch,
		queue:        q,
		emailSender:  NewSMTPEmailSender(),
		inAppSender:  NewDBInAppSender(queries),
	}, nil
}

func (s *NotificationService) SendNotification(ctx context.Context, notification Notification) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	// Save notification to database
	_, err = qtx.CreateNotification(ctx, models.CreateNotificationParams{
		Userid:  0,
		Type:    string(notification.Type),
		Message: notification.Content,
	})
	if err != nil {
		return fmt.Errorf("failed to create notification in database: %v", err)
	}

	// Publish to RabbitMQ
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}

	err = s.channel.PublishWithContext(ctx,
		"",              // exchange
		s.queue.Name,    // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (s *NotificationService) StartConsumer(ctx context.Context) error {
	msgs, err := s.channel.Consume(
		s.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			var notification Notification
			if err := json.Unmarshal(d.Body, &notification); err != nil {
				log.Printf("Error unmarshalling notification: %v", err)
				d.Nack(false, true)
				continue
			}

			if err := s.handleNotification(notification); err != nil {
				log.Printf("Error handling notification: %v", err)
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (s *NotificationService) handleNotification(notification Notification) error {
	switch notification.Type {
	case EmailNotification:
		return s.emailSender.SendEmail(notification.To, notification.Subject, notification.Content)
	case InAppNotification:
		return s.inAppSender.SendInAppNotification(notification.To, notification.Content)
	default:
		return fmt.Errorf("unknown notification type: %s", notification.Type)
	}
}

func (s *NotificationService) Close() error {
	if err := s.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %v", err)
	}
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}
	return nil
}

func (s *NotificationService) GetUnreadNotifications(ctx context.Context, userID int32) ([]models.Notification, error) {
	return s.queries.GetUnreadNotifications(ctx, userID)
}

func (s *NotificationService) MarkNotificationsAsRead(ctx context.Context, userID int32) error {
	return s.queries.MarkNotificationsAsRead(ctx, userID)
}