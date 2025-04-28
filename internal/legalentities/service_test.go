package legalentities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/krisch/crm-backend/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAll() ([]domain.LegalEntity, error) {
	args := m.Called()
	return args.Get(0).([]domain.LegalEntity), args.Error(1)
}

func (m *MockRepository) Create(entity *domain.LegalEntity) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockRepository) Update(entity *domain.LegalEntity) error {
	args := m.Called(entity)
	return args.Error(0)
}

func (m *MockRepository) Delete(uuid string) error {
	args := m.Called(uuid)
	return args.Error(0)
}

func (m *MockRepository) GetAllBankAccounts(legalEntityUUID uuid.UUID) ([]domain.BankAccount, error) {
	args := m.Called(legalEntityUUID)
	return args.Get(0).([]domain.BankAccount), args.Error(1)
}

func (m *MockRepository) CreateBankAccount(account *domain.BankAccount) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockRepository) UpdateBankAccount(account *domain.BankAccount) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockRepository) DeleteBankAccount(accountUUID uuid.UUID) error {
	args := m.Called(accountUUID)
	return args.Error(0)
}

func (m *MockRepository) GetBankAccount(accountUUID uuid.UUID) (*domain.BankAccount, error) {
	args := m.Called(accountUUID)
	return args.Get(0).(*domain.BankAccount), args.Error(1)
}

func (m *MockRepository) ClearPrimaryFlag(legalEntityID uuid.UUID) error {
	args := m.Called(legalEntityID)
	return args.Error(0)
}

func TestCreateBankAccount(t *testing.T) {
	legalEntityID := uuid.New()
	accountID := uuid.New()

	tests := []struct {
		name        string
		account     *domain.BankAccount
		setupMock   func(repo *MockRepository)
		expectError bool
	}{
		{
			name: "successful creation",
			account: &domain.BankAccount{
				UUID:                 accountID,
				LegalEntityID:        legalEntityID,
				BIC:                  "044525974",
				BankName:             "Test Bank",
				SettlementAccount:    "40702810716540010359",
				CorrespondentAccount: "30101810400000000225",
				Currency:             "RUB",
				IsPrimary:            true,
			},
			setupMock: func(repo *MockRepository) {
				repo.On("ClearPrimaryFlag", legalEntityID).Return(nil)
				repo.On("CreateBankAccount", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "invalid BIK",
			account: &domain.BankAccount{
				BIC:               "123",
				SettlementAccount: "40702810716540010359",
			},
			setupMock:   func(repo *MockRepository) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{}
			tt.setupMock(repo)

			service := &ServiceImpl{repo: repo}
			err := service.CreateBankAccount(tt.account)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				repo.AssertExpectations(t)
			}
		})
	}
}

func TestGetAllBankAccounts(t *testing.T) {
	legalEntityID := uuid.New()
	accountID := uuid.New()

	tests := []struct {
		name          string
		legalEntityID uuid.UUID
		setupMock     func(repo *MockRepository)
		expectLen     int
		expectError   bool
	}{
		{
			name:          "successful get all",
			legalEntityID: legalEntityID,
			setupMock: func(repo *MockRepository) {
				repo.On("GetAllBankAccounts", legalEntityID).Return([]domain.BankAccount{
					{
						UUID:              accountID,
						LegalEntityID:     legalEntityID,
						BIC:               "044525974",
						BankName:          "Test Bank",
						SettlementAccount: "40702810716540010359",
					},
				}, nil)
			},
			expectLen:   1,
			expectError: false,
		},
		{
			name:          "not found",
			legalEntityID: legalEntityID,
			setupMock: func(repo *MockRepository) {
				repo.On("GetAllBankAccounts", legalEntityID).Return([]domain.BankAccount{}, nil)
			},
			expectLen:   0,
			expectError: false,
		},
		{
			name:          "repository error",
			legalEntityID: legalEntityID,
			setupMock: func(repo *MockRepository) {
				repo.On("GetAllBankAccounts", legalEntityID).Return(nil, assert.AnError)
			},
			expectLen:   0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{}
			tt.setupMock(repo)

			service := &ServiceImpl{repo: repo}
			accounts, err := service.GetAllBankAccounts(tt.legalEntityID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, accounts, tt.expectLen)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestUpdateBankAccount(t *testing.T) {
	legalEntityID := uuid.New()
	accountID := uuid.New()

	tests := []struct {
		name        string
		account     *domain.BankAccount
		setupMock   func(repo *MockRepository)
		expectError bool
	}{
		{
			name: "successful update",
			account: &domain.BankAccount{
				UUID:              accountID,
				LegalEntityID:     legalEntityID,
				BIC:               "044525974",
				BankName:          "Updated Bank",
				SettlementAccount: "40702810716540010359",
				IsPrimary:         true,
			},
			setupMock: func(repo *MockRepository) {
				repo.On("GetBankAccount", accountID).Return(&domain.BankAccount{
					UUID:          accountID,
					LegalEntityID: legalEntityID,
					IsPrimary:     false,
				}, nil)
				repo.On("ClearPrimaryFlag", legalEntityID).Return(nil)
				repo.On("UpdateBankAccount", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "account not found",
			account: &domain.BankAccount{
				UUID: accountID,
			},
			setupMock: func(repo *MockRepository) {
				repo.On("GetBankAccount", accountID).Return(nil, ErrBankAccountNotFound)
			},
			expectError: true,
		},
		{
			name: "invalid data",
			account: &domain.BankAccount{
				UUID:              accountID,
				LegalEntityID:     legalEntityID,
				BIC:               "123", // invalid
				SettlementAccount: "40702810716540010359",
			},
			setupMock: func(repo *MockRepository) {
				repo.On("GetBankAccount", accountID).Return(&domain.BankAccount{
					UUID: accountID,
				}, nil)
			},
			expectError: true,
		},
		{
			name: "clear primary error",
			account: &domain.BankAccount{
				UUID:              accountID,
				LegalEntityID:     legalEntityID,
				BIC:               "044525974",
				BankName:          "Bank",
				SettlementAccount: "40702810716540010359",
				IsPrimary:         true,
			},
			setupMock: func(repo *MockRepository) {
				repo.On("GetBankAccount", accountID).Return(&domain.BankAccount{
					UUID:          accountID,
					LegalEntityID: legalEntityID,
					IsPrimary:     false,
				}, nil)
				repo.On("ClearPrimaryFlag", legalEntityID).Return(assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{}
			tt.setupMock(repo)

			service := &ServiceImpl{repo: repo}
			err := service.UpdateBankAccount(tt.account)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestDeleteBankAccount(t *testing.T) {
	accountID := uuid.New()

	tests := []struct {
		name        string
		accountUUID uuid.UUID
		setupMock   func(repo *MockRepository)
		expectError bool
	}{
		{
			name:        "successful delete",
			accountUUID: accountID,
			setupMock: func(repo *MockRepository) {
				repo.On("DeleteBankAccount", accountID).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "repository error",
			accountUUID: accountID,
			setupMock: func(repo *MockRepository) {
				repo.On("DeleteBankAccount", accountID).Return(assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{}
			tt.setupMock(repo)

			service := &ServiceImpl{repo: repo}
			err := service.DeleteBankAccount(tt.accountUUID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestGetBankAccount(t *testing.T) {
	accountID := uuid.New()
	legalEntityID := uuid.New()

	tests := []struct {
		name        string
		accountUUID uuid.UUID
		setupMock   func(repo *MockRepository)
		expectAcc   *domain.BankAccount
		expectError bool
	}{
		{
			name:        "successful get",
			accountUUID: accountID,
			setupMock: func(repo *MockRepository) {
				repo.On("GetBankAccount", accountID).Return(&domain.BankAccount{
					UUID:              accountID,
					LegalEntityID:     legalEntityID,
					BIC:               "044525974",
					BankName:          "Test Bank",
					SettlementAccount: "40702810716540010359",
				}, nil)
			},
			expectAcc: &domain.BankAccount{
				UUID:              accountID,
				LegalEntityID:     legalEntityID,
				BIC:               "044525974",
				BankName:          "Test Bank",
				SettlementAccount: "40702810716540010359",
			},
			expectError: false,
		},
		{
			name:        "not found",
			accountUUID: accountID,
			setupMock: func(repo *MockRepository) {
				repo.On("GetBankAccount", accountID).Return(nil, ErrBankAccountNotFound)
			},
			expectAcc:   nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{}
			tt.setupMock(repo)

			service := &ServiceImpl{repo: repo}
			acc, err := service.GetBankAccount(tt.accountUUID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectAcc, acc)
			}
			repo.AssertExpectations(t)
		})
	}
}
