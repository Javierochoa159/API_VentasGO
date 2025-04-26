package metadata

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

// Create adds a brand-new metadata to the system.
// It sets CreatedAt and UpdatedAt to the current time and initializes Version to 1.
// Returns ErrEmptyID if Metadata.ID is empty.
func (s *Service) Create(metadata *Metadata, userId string) error {
	metadata.Approved = 0
	metadata.Pending = 0
	metadata.Quantity = 0
	metadata.Rejected = 0
	metadata.Total_amount = 0.0

	err := s.storage.SetMetadata(metadata, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(status string, id string, amount ...float32) (*Metadata, error) {
	existing, err := s.storage.ReadMetadata(id)
	if err != nil {
		return nil, err
	}

	switch status {
	case "approved":
		s.changeToApproved(id)
	case "rejected":
		s.changeToRejected(id)
	case "new_sale":
		if len(amount) == 0 {
			return nil, ErrNotValidOperation
		}
		if amount[0] == 0 {
			return nil, ErrNotValidOperation
		}
		s.incrementSale(id, amount[0])
	default:
		return nil, ErrNotValidOperation
	}

	return existing, nil
}

func (s *Service) Get(userId string) *Metadata {
	return s.storage.mapMeta[userId]
}

func (s *Service) incrementSale(userId string, totalAmount float32) {
	s.storage.mapMeta[userId].Quantity++
	s.storage.mapMeta[userId].Pending++
	s.storage.mapMeta[userId].Total_amount += totalAmount
}

func (s *Service) changeToApproved(userId string) {
	s.storage.mapMeta[userId].Approved++
	s.storage.mapMeta[userId].Pending--
}

func (s *Service) changeToRejected(userId string) {
	s.storage.mapMeta[userId].Rejected++
	s.storage.mapMeta[userId].Pending--
}
