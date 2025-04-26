package metadata

import "errors"

// ErrEmpty is returned when trying to store a metada with an empty ID.
var ErrEmptyIDMetadata = errors.New("empty metadata")

// ErrNotFound is returned when a metadata with the given ID is not found.
var ErrNotFoundMetadata = errors.New("metadata not found")

// ErrNo inValidOperation is returned when the user performs an invalid operation.
var ErrNotValidOperation = errors.New("invalid operation")

// LocalStorage provides an in-memory implementation for storing sales.
type LocalStorage struct {
	mapMeta map[string]*Metadata
}

// NewLocalStorage instantiates a new LocalStorage with an empty map.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		mapMeta: map[string]*Metadata{},
	}
}

// Set stores or updates a metadata in the local storage.
// Returns ErrEmptyIDMetadata if the metadata has an empty ID.
func (l *LocalStorage) SetMetadata(metadata *Metadata, id string) error {
	if id == "" {
		return ErrEmptyIDMetadata
	}
	l.mapMeta[id] = metadata
	return nil
}

// Read retrieves a metadata from the local storage by ID.
// Returns ErrNotFoundMetadata if the metadata is not found.
func (l *LocalStorage) ReadMetadata(id string) (*Metadata, error) {
	u, ok := l.mapMeta[id]
	if !ok {
		return nil, ErrNotFoundMetadata
	}

	return u, nil
}

// Delete removes metadata from the local storage by ID.
// Returns ErrNotFoundMetadata if the metadata does not exist.
func (l *LocalStorage) DeleteMetadata(id string) error {
	_, err := l.ReadMetadata(id)
	if err != nil {
		return err
	}

	delete(l.mapMeta, id)
	return nil
}
