package senders

import (
	"context"
	"log"
	"sync"

	"github.com/iRankHub/backend/internal/services/notification/models"
)

type InAppSender struct {
	subscribers map[string][]chan *models.Notification
	mu          sync.RWMutex
}

func NewInAppSender() *InAppSender {
	return &InAppSender{
		subscribers: make(map[string][]chan *models.Notification),
	}
}

func (s *InAppSender) Send(ctx context.Context, notification *models.Notification) error {
	s.mu.RLock()
	channels, exists := s.subscribers[notification.UserID]
	s.mu.RUnlock()

	if !exists || len(channels) == 0 {
		// No active subscribers, mark as pending
		notification.UpdateDeliveryStatus(models.InAppDelivery, models.StatusPending, nil)
		return nil
	}

	// Fan out notification to all subscribers
	for _, ch := range channels {
		select {
		case ch <- notification:
			notification.UpdateDeliveryStatus(models.InAppDelivery, models.StatusDelivered, nil)
		default:
			// Channel is full, skip this subscriber
			log.Printf("Subscriber channel full for user %s", notification.UserID)
		}
	}

	return nil
}

func (s *InAppSender) Subscribe(ctx context.Context, userID string) (<-chan *models.Notification, func(), error) {
	ch := make(chan *models.Notification, 100) // Buffer size of 100

	s.mu.Lock()
	if s.subscribers[userID] == nil {
		s.subscribers[userID] = make([]chan *models.Notification, 0)
	}
	s.subscribers[userID] = append(s.subscribers[userID], ch)
	s.mu.Unlock()

	// Create unsubscribe function
	unsubscribe := func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		channels := s.subscribers[userID]
		for i, c := range channels {
			if c == ch {
				// Remove this channel
				s.subscribers[userID] = append(channels[:i], channels[i+1:]...)
				break
			}
		}
		close(ch)

		// If no more subscribers for this user, clean up
		if len(s.subscribers[userID]) == 0 {
			delete(s.subscribers, userID)
		}
	}

	// Start cleanup goroutine
	go func() {
		<-ctx.Done()
		unsubscribe()
	}()

	return ch, unsubscribe, nil
}

func (s *InAppSender) UnsubscribeAll(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	channels, exists := s.subscribers[userID]
	if !exists {
		return
	}

	// Close all channels
	for _, ch := range channels {
		close(ch)
	}

	delete(s.subscribers, userID)
}

func (s *InAppSender) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Close all subscriber channels
	for userID, channels := range s.subscribers {
		for _, ch := range channels {
			close(ch)
		}
		delete(s.subscribers, userID)
	}

	return nil
}
