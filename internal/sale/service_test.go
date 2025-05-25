package sale

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	type fields struct {
		storage     Storage
		userService UserService
	}

	type args struct {
		sale *Sale
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  func(t *testing.T, err error)
		wantSale func(t *testing.T, sale *Sale)
	}{
		{
			name: "error",
			fields: fields{
				storage: &mockStorageSale{
					mockSetSale: func(sale *Sale) error {
						return errors.New("fake error trying to set sale")
					},
				},
			},
			args: args{
				sale: &Sale{},
			},
			wantErr: func(t *testing.T, err error) {
				require.NotNil(t, err)
				require.EqualError(t, err, "fake error trying to set sale")
			},
			wantSale: nil,
		},
		{
			name: "errorUserNotFound",
			fields: fields{
				storage: NewLocalStorage(),
				userService: &mockUserService{
					mockFindUser: func(id string) error {
						return errors.New("user not found")
					},
				},
			},
			args: args{
				sale: &Sale{
					UserId: "1000",
					Amount: 1500,
				},
			},
			wantErr: func(t *testing.T, err error) {
				require.NotNil(t, err)
				require.EqualError(t, err, "user not found")
			},
			wantSale: nil,
		},
		{
			name: "success",
			fields: fields{
				storage: NewLocalStorage(),
			},
			args: args{
				sale: &Sale{
					UserId: "1",
					Amount: 1500,
				},
			},
			wantErr: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
			wantSale: func(t *testing.T, input *Sale) {
				require.NotEmpty(t, input.ID)
				require.NotEmpty(t, input.Amount)
				require.NotEmpty(t, input.Status)
				require.NotEmpty(t, input.CreatedAt)
				require.NotEmpty(t, input.UpdatedAt)
				require.Equal(t, 1, input.Version)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.storage, tt.fields.userService, nil)

			err := s.Create(tt.args.sale)
			if tt.wantErr != nil {
				tt.wantErr(t, err)
			}

			if tt.wantSale != nil {
				tt.wantSale(t, tt.args.sale)
			}
		})
	}
}

type mockStorageSale struct {
	mockSetSale                  func(sale *Sale) error
	mockReadSale                 func(id string) (*Sale, error)
	mockDeleteSale               func(id string) error
	mockReadSalesByUser          func(id string) ([]*Sale, map[string]float32)
	mockReadSalesByUserAndStatus func(id string, status string) ([]*Sale, map[string]float32)
}

func (m *mockStorageSale) SetSale(sale *Sale) error {
	return m.mockSetSale(sale)
}

func (m *mockStorageSale) ReadSale(id string) (*Sale, error) {
	return m.mockReadSale(id)
}

func (m *mockStorageSale) ReadSalesByUser(id string) ([]*Sale, map[string]float32) {
	return m.mockReadSalesByUser(id)
}

func (m *mockStorageSale) ReadSalesByUserAndStatus(id string, status string) ([]*Sale, map[string]float32) {
	return m.mockReadSalesByUserAndStatus(id, status)
}

func (m *mockStorageSale) DeleteSale(id string) error {
	return m.mockDeleteSale(id)
}

type mockUserService struct {
	mockFindUser func(id string) error
}

func (m *mockUserService) FindUser(id string) error {
	if m.mockFindUser != nil {
		return m.mockFindUser(id)
	}
	return nil
}
