package sale

import (
	"time"

	"github.com/google/uuid"
)

// Service provides high-level sale management operations on a LocalStorage backend.
type Service struct {
	// storage is the underlying persistence for Sale entities.
	storage *LocalStorage
}

// NewService creates a new Service.
func NewService(storage *LocalStorage) *Service {
	return &Service{
		storage: storage,
	}
}

// Create adds a brand-new sale to the system.
// It sets CreatedAt and UpdatedAt to the current time and initializes Version to 1.
// Returns ErrEmptyID if sale.ID is empty.
func (s *Service) Create(sale *Sale) error {
	sale.ID = uuid.NewString()
	now := time.Now()
	sale.CreatedAt = now
	sale.UpdatedAt = now
	sale.Version = 1

	return s.storage.Set(sale)
}

// Get retrieves a sale by its ID.
// Returns ErrNotFound if no user exists with the given ID.
func (s *Service) Get(id string) (*Sale, error) {
	return s.storage.Read(id)
}

// Update modifies an existing sale's data.
// It updates Status, sets UpdatedAt to now and increments Version.
// Returns ErrNotFound if the sale does not exist, or ErrEmptyID if sale.ID is empty.
func (s *Service) Update(id string, sale *UpdateFields) (*Sale, error) {
	existing, err := s.storage.Read(id)
	if err != nil {
		return nil, err
	}

	if sale.Status != nil {
		existing.Status = *sale.Status
	}

	existing.UpdatedAt = time.Now()
	existing.Version++

	if err := s.storage.Set(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete removes a sale from the system by its ID.
// Returns ErrNotFound if the sale does not exist.
func (s *Service) Delete(id string) error {
	return s.storage.Delete(id)
}
