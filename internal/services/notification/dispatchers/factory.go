package dispatchers

import (
	"fmt"
	"sync"

	"github.com/iRankHub/backend/internal/services/notification/models"
)

// DispatcherFactory manages the creation and caching of dispatchers
type DispatcherFactory struct {
	baseDispatcher *BaseDispatcher
	options        DispatcherOptions
	dispatchers    map[models.Category]Dispatcher
	mu             sync.RWMutex
}

// NewDispatcherFactory creates a new dispatcher factory
func NewDispatcherFactory(base *BaseDispatcher, options DispatcherOptions) *DispatcherFactory {
	return &DispatcherFactory{
		baseDispatcher: base,
		options:        options,
		dispatchers:    make(map[models.Category]Dispatcher),
	}
}

// GetDispatcher returns the appropriate dispatcher for a notification category
func (f *DispatcherFactory) GetDispatcher(category models.Category) (Dispatcher, error) {
	f.mu.RLock()
	if dispatcher, exists := f.dispatchers[category]; exists {
		f.mu.RUnlock()
		return dispatcher, nil
	}
	f.mu.RUnlock()

	// Create new dispatcher if it doesn't exist
	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check after acquiring write lock
	if dispatcher, exists := f.dispatchers[category]; exists {
		return dispatcher, nil
	}

	var dispatcher Dispatcher
	switch category {
	case models.AuthCategory:
		dispatcher = NewAuthDispatcher(f.baseDispatcher, f.options)
	case models.UserCategory:
		dispatcher = NewUserDispatcher(f.baseDispatcher, f.options)
	case models.TournamentCategory:
		dispatcher = NewTournamentDispatcher(f.baseDispatcher, f.options)
	case models.DebateCategory:
		dispatcher = NewDebateDispatcher(f.baseDispatcher, f.options)
	case models.ReportCategory:
		dispatcher = NewReportDispatcher(f.baseDispatcher, f.options)
	default:
		return nil, fmt.Errorf("unknown notification category: %s", category)
	}

	f.dispatchers[category] = dispatcher
	return dispatcher, nil
}

// UpdateOptions updates the dispatcher options
func (f *DispatcherFactory) UpdateOptions(options DispatcherOptions) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.options = options
	// Clear cache to force recreation with new options
	f.dispatchers = make(map[models.Category]Dispatcher)
}
