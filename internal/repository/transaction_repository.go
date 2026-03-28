package repository

import (
	"time"

	"github.com/bakdaaswandi5818/pencatat-keuangan/internal/domain"
	"github.com/google/uuid"
)

// Filter holds optional query parameters for listing transactions.
type Filter struct {
	Type     string
	Category string
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
}

// TransactionRepository defines the persistence interface.
// Using an interface here allows seamless migration to PostgreSQL in the future.
type TransactionRepository interface {
	Create(tx *domain.Transaction) error
	GetByID(id uuid.UUID) (*domain.Transaction, error)
	List(f Filter) ([]domain.Transaction, int64, error)
	Update(tx *domain.Transaction) error
	Delete(id uuid.UUID) error
	Summary() (*domain.Summary, error)
}
