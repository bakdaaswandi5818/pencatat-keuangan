package service

import (
	"errors"
	"time"

	"github.com/bakdaaswandi5818/pencatat-keuangan/internal/domain"
	"github.com/bakdaaswandi5818/pencatat-keuangan/internal/repository"
	"github.com/google/uuid"
)

// CreateTransactionInput holds the data required to create a transaction.
type CreateTransactionInput struct {
	Title           string                `json:"title"            validate:"required"`
	Amount          float64               `json:"amount"           validate:"required,gt=0"`
	Type            domain.TransactionType `json:"type"             validate:"required,oneof=income expense"`
	Category        string                `json:"category"         validate:"required"`
	TransactionDate time.Time             `json:"transaction_date" validate:"required"`
}

// ListTransactionsInput holds query parameters for listing transactions.
type ListTransactionsInput struct {
	Type     string
	Category string
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
}

// ListTransactionsOutput is returned by ListTransactions.
type ListTransactionsOutput struct {
	Data  []domain.Transaction `json:"data"`
	Total int64                `json:"total"`
	Limit int                  `json:"limit"`
	Offset int                 `json:"offset"`
}

// TransactionService defines the business-logic interface.
type TransactionService interface {
	Create(input CreateTransactionInput) (*domain.Transaction, error)
	GetByID(id uuid.UUID) (*domain.Transaction, error)
	List(input ListTransactionsInput) (*ListTransactionsOutput, error)
	Delete(id uuid.UUID) error
	GetSummary() (*domain.Summary, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) Create(input CreateTransactionInput) (*domain.Transaction, error) {
	tx := &domain.Transaction{
		Title:           input.Title,
		Amount:          input.Amount,
		Type:            input.Type,
		Category:        input.Category,
		TransactionDate: input.TransactionDate,
	}
	if err := s.repo.Create(tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *transactionService) GetByID(id uuid.UUID) (*domain.Transaction, error) {
	tx, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *transactionService) List(input ListTransactionsInput) (*ListTransactionsOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	filter := repository.Filter{
		Type:     input.Type,
		Category: input.Category,
		DateFrom: input.DateFrom,
		DateTo:   input.DateTo,
		Limit:    limit,
		Offset:   input.Offset,
	}

	txs, total, err := s.repo.List(filter)
	if err != nil {
		return nil, err
	}

	return &ListTransactionsOutput{
		Data:   txs,
		Total:  total,
		Limit:  limit,
		Offset: input.Offset,
	}, nil
}

func (s *transactionService) Delete(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}
	return s.repo.Delete(id)
}

func (s *transactionService) GetSummary() (*domain.Summary, error) {
	return s.repo.Summary()
}
