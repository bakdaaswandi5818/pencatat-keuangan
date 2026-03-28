package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionType defines allowed transaction types.
type TransactionType string

const (
	TypeIncome  TransactionType = "income"
	TypeExpense TransactionType = "expense"
)

// Transaction is the core domain model.
type Transaction struct {
	ID              uuid.UUID       `gorm:"type:text;primaryKey" json:"id"`
	Title           string          `gorm:"not null"             json:"title"`
	Amount          float64         `gorm:"not null"             json:"amount"`
	Type            TransactionType `gorm:"index;not null"       json:"type"`
	Category        string          `gorm:"index"                json:"category"`
	TransactionDate time.Time       `gorm:"index;not null"       json:"transaction_date"`
	CreatedAt       time.Time       `                            json:"created_at"`
	UpdatedAt       time.Time       `                            json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index"                json:"deleted_at,omitempty"`
}

// BeforeCreate sets a new UUID before inserting a record.
func (t *Transaction) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// Summary holds aggregated financial totals.
type Summary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}
