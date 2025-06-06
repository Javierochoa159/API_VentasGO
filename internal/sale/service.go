package sale

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService interface {
	FindUser(id string) error
}

// Service provides high-level sale management operations on a LocalStorage backend.
type Service struct {
	// storage is the underlying persistence for User entities.
	storage Storage

	// logger is our observability component to log.
	Logger *zap.Logger

	userService UserService
}

// NewService creates a new Service.
func NewService(storage Storage, userService UserService, logger *zap.Logger) *Service {
	if logger == nil {
		logger, _ = zap.NewProduction()
		defer logger.Sync()
	}

	return &Service{
		storage:     storage,
		userService: userService,
		Logger:      logger,
	}
}

// Create adds a brand-new sale to the system.
// It sets CreatedAt and UpdatedAt to the current time and initializes Version to 1.
// Returns ErrEmptyID if sale.ID is empty.
func (s *Service) Create(sale *Sale) error {
	if s.userService != nil {
		if err := s.userService.FindUser(sale.UserId); err != nil {
			s.Logger.Error("user not found", zap.Error(err))
			return err
		}
	}

	if os.Getenv("MODO") != "testing" {
		opciones := []string{"approved", "rejected", "pending"}
		estado := rand.Intn(3)
		sale.Status = opciones[estado]
	} else {
		sale.Status = "pending"
	}

	sale.ID = uuid.NewString()
	now := time.Now()
	sale.CreatedAt = now
	sale.UpdatedAt = now
	sale.Version = 1
	if err := s.storage.SetSale(sale); err != nil {
		s.Logger.Error("failed to set sale", zap.Error(err), zap.Any("sale", sale))
		return err
	}

	return nil
}

// Get retrieves a sale by its ID.
// Returns ErrNotFound if no user exists with the given ID.
func (s *Service) Get(id string) (*Sale, error) {
	return s.storage.ReadSale(id)
}

func (s *Service) GetUserSales(id string, status string) ([]*Sale, map[string]float32) {
	if status == "" {
		return s.storage.ReadSalesByUser(id)
	}
	status = strings.ToLower(status)
	return s.storage.ReadSalesByUserAndStatus(id, status)
}

// Update modifies an existing sale's data.
// It updates Status, sets UpdatedAt to now and increments Version.
// Returns ErrNotFound if the sale does not exist, or ErrEmptyID if sale.ID is empty.
// Returns ErrNotValidOperation if the sale status is invalid for the operation.
func (s *Service) Update(id string, sale *UpdateFields) (*Sale, error) {
	existing, err := s.storage.ReadSale(id)
	if err != nil {
		return nil, err

	}

	if existing.Status != "pending" {
		return nil, ErrInvalidStatus
	}

	if sale.Status != nil {
		existing.Status = strings.ToLower(*sale.Status)
	}

	existing.UpdatedAt = time.Now()
	existing.Version++

	if err := s.storage.SetSale(existing); err != nil {
		return nil, err
	}

	return existing, nil
}
