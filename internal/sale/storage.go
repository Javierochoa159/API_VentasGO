package sale

import (
	"errors"
)

// ErrNotFound is returned when a sale with the given ID is not found.
var ErrNotFound = errors.New("sale not found")

// ErrStatusNotFound is returned when a status not found.
var ErrStatusNotFound = errors.New("sale status not found")

// ErrEmptyID is returned when trying to store a sale with an empty ID.
var ErrEmptyID = errors.New("empty sale ID")

// ErrInvalidAmoun is returned when trying to store a sale with an empty amount.
var ErrInvalidAmoun = errors.New("amount equals or lower 0")

// ErrInvalidStatus is returned when the user performs an invalid status.
var ErrInvalidStatus = errors.New("invalid status")

// ErrNo inValidOperation is returned when the user performs an invalid operation.
var ErrNotValidOperation = errors.New("invalid operation")

type Storage interface {
	SetSale(sale *Sale) error
	ReadSale(id string) (*Sale, error)
	ReadSalesByUser(id string) ([]*Sale, map[string]float32)
	ReadSalesByUserAndStatus(id string, status string) ([]*Sale, map[string]float32)
	DeleteSale(id string) error
}

// LocalStorage provides an in-memory implementation for storing sales.
type LocalStorage struct {
	mapSale map[string]*Sale
}

// NewLocalStorage instantiates a new LocalStorage with an empty map.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		mapSale: make(map[string]*Sale),
	}
}

// Set stores or updates a sale in the local storage.
// Returns ErrEmptyID if the sale has an empty ID.
func (l *LocalStorage) SetSale(sale *Sale) error {
	if sale.ID == "" {
		return ErrEmptyID
	}

	l.mapSale[sale.ID] = sale
	return nil
}

// Read retrieves a sale from the local storage by ID.
// Returns ErrNotFound if the sale is not found.
func (l *LocalStorage) ReadSale(id string) (*Sale, error) {
	u, ok := l.mapSale[id]
	if !ok {
		return nil, ErrNotFound
	}

	return u, nil
}

func (l *LocalStorage) ReadSalesByUser(id string) ([]*Sale, map[string]float32) {
	meta := map[string]float32{
		"quantity":     0,
		"approved":     0,
		"pending":      0,
		"rejected":     0,
		"total_amount": 0,
	}
	var sales []*Sale
	for _, sale := range l.mapSale {
		if sale.UserId == id {
			sales = append(sales, sale)
			meta["quantity"]++
			meta["total_amount"] += sale.Amount
			switch sale.Status {
			case "approved":
				meta["approved"]++
			case "rejected":
				meta["rejected"]++
			case "pending":
				meta["pending"]++
			}
		}
	}
	return sales, meta
}

func (l *LocalStorage) ReadSalesByUserAndStatus(id string, status string) ([]*Sale, map[string]float32) {
	meta := map[string]float32{
		"quantity":     0,
		"approved":     0,
		"pending":      0,
		"rejected":     0,
		"total_amount": 0,
	}
	var sales []*Sale
	for _, sale := range l.mapSale {
		if sale.UserId == id && sale.Status == status {
			sales = append(sales, sale)
			meta["quantity"]++
			meta["total_amount"] += sale.Amount
			switch sale.Status {
			case "approved":
				meta["approved"]++
			case "rejected":
				meta["rejected"]++
			case "pending":
				meta["pending"]++
			}
		}
	}
	return sales, meta
}

// Delete removes a sale from the local storage by ID.
// Returns ErrNotFound if the sale does not exist.
func (l *LocalStorage) DeleteSale(id string) error {
	_, err := l.ReadSale(id)
	if err != nil {
		return err
	}

	delete(l.mapSale, id)
	return nil
}
