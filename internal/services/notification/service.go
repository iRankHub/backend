package notification

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	pb "github.com/iRankHub/backend/internal/grpc/proto/notification"
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
	db            *sql.DB
	queries       *models.Queries
	conn          *amqp.Connection
	channel       *amqp.Channel
	queue         amqp.Queue
	emailSender   EmailSender
	inAppSender   InAppSender
	subscriptions map[int32][]chan *pb.Notification
	mu            sync.RWMutex
}

func NewNotificationService(db *sql.DB) (*NotificationService, error) {
	queries := models.New(db)

	emailSender, err := NewSMTPEmailSender()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize email sender: %v", err)
	}

	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return nil, fmt.Errorf("RABBITMQ_URL environment variable is not set")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	// Set QoS for better message handling
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %v", err)
	}

	q, err := ch.QueueDeclare(
		"notifications",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	service := &NotificationService{
		db:            db,
		queries:       queries,
		conn:          conn,
		channel:       ch,
		queue:         q,
		emailSender:   emailSender,
		inAppSender:   NewDBInAppSender(queries),
		subscriptions: make(map[int32][]chan *pb.Notification),
	}

	go service.startConsumer()

	return service, nil
}

func (s *NotificationService) startConsumer() {
	backoff := time.Second
	maxBackoff := 30 * time.Second

	for {
		msgs, err := s.channel.Consume(
			s.queue.Name,
			"",    // consumer
			false, // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			log.Printf("Failed to start consuming messages: %v. Retrying in %v...", err, backoff)
			time.Sleep(backoff)
			backoff = min(backoff*2, maxBackoff)
			continue
		}

		log.Println("Started consuming messages")
		backoff = time.Second // Reset backoff on successful connection

		for msg := range msgs {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			var notification Notification
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				log.Printf("Error unmarshalling notification: %v", err)
				msg.Nack(false, false) // Don't requeue on unmarshal error
				cancel()
				continue
			}

			err := s.handleNotificationWithTimeout(ctx, notification)
			if err != nil {
				log.Printf("Error handling notification for %s: %v", notification.To, err)
				shouldRequeue := !isPermanentError(err)
				msg.Nack(false, shouldRequeue)
			} else {
				msg.Ack(false)
			}
			cancel()
		}

		log.Println("Consumer channel closed. Attempting to reconnect...")
		time.Sleep(backoff)
	}
}

func (s *NotificationService) handleNotificationWithTimeout(ctx context.Context, notification Notification) error {
	done := make(chan error, 1)

	go func() {
		done <- s.handleNotification(notification)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("notification processing timed out")
	}
}

func isPermanentError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Define permanent error conditions
	permanentErrors := []string{
		"unknown notification type",
		"invalid email address",
		"recipient not found",
		"failed to unmarshal",
		"permanent error",
	}

	for _, pe := range permanentErrors {
		if strings.Contains(errStr, pe) {
			return true
		}
	}

	// Email-specific permanent errors
	emailPermanentErrors := []string{
		"invalid auth",
		"invalid recipient",
		"malformed email",
		"user unknown",
		"email address rejected",
	}

	for _, epe := range emailPermanentErrors {
		if strings.Contains(errStr, epe) {
			return true
		}
	}

	return false
}

func (s *NotificationService) handleNotification(notification Notification) error {
	switch notification.Type {
	case EmailNotification:
		if !isValidEmail(notification.To) {
			return fmt.Errorf("permanent error: invalid email address: %s", notification.To)
		}
		if err := s.emailSender.SendEmail(notification.To, notification.Subject, notification.Content); err != nil {
			return fmt.Errorf("email sending failed: %w", err)
		}
	case InAppNotification:
		if notification.To == "" {
			return fmt.Errorf("permanent error: recipient not found")
		}
		if err := s.inAppSender.SendInAppNotification(notification.To, notification.Content); err != nil {
			return fmt.Errorf("in-app notification failed: %w", err)
		}
	default:
		return fmt.Errorf("permanent error: unknown notification type: %s", notification.Type)
	}
	return nil
}

// Helper function to validate email addresses
func isValidEmail(email string) bool {
	// Basic email validation
	if email == "" {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

func (s *NotificationService) RegisterNotificationChannel(userID int32, ch chan *pb.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscriptions[userID] = append(s.subscriptions[userID], ch)
}

func (s *NotificationService) UnregisterNotificationChannel(userID int32, ch chan *pb.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	channels := s.subscriptions[userID]
	for i, c := range channels {
		if c == ch {
			s.subscriptions[userID] = append(channels[:i], channels[i+1:]...)
			close(ch)
			break
		}
	}
}

func (s *NotificationService) SendNotification(ctx context.Context, notification Notification) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	// Save notification to database
	userID, _ := strconv.Atoi(notification.To)
	createdNotification, err := qtx.CreateNotification(ctx, models.CreateNotificationParams{
		Userid:         int32(userID),
		Type:           string(notification.Type),
		Recipientemail: sql.NullString{String: notification.To, Valid: notification.To != ""},
		Subject:        sql.NullString{String: notification.Subject, Valid: notification.Subject != ""},
		Message:        notification.Content,
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
		"",           // exchange
		s.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
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

	// Send to subscribed channels
	s.mu.RLock()
	channels := s.subscriptions[int32(userID)]
	s.mu.RUnlock()

	protoNotification := &pb.Notification{
		Id:      int32(createdNotification.Notificationid),
		Type:    pb.NotificationType(pb.NotificationType_value[string(notification.Type)]),
		To:      notification.To,
		Subject: notification.Subject,
		Content: notification.Content,
	}

	for _, ch := range channels {
		select {
		case ch <- protoNotification:
		default:
			log.Printf("Failed to send notification to channel for user %d: channel full or closed", userID)
		}
	}

	return nil
}

func (s *NotificationService) SubscribeToNotifications(userID int32) (<-chan *pb.Notification, func()) {
	ch := make(chan *pb.Notification, 100)
	s.RegisterNotificationChannel(userID, ch)

	return ch, func() {
		s.UnregisterNotificationChannel(userID, ch)
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

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
