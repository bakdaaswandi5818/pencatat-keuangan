package repository

import (
	"time"

	"github.com/bakdaaswandi5818/pencatat-keuangan/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormTransactionRepository struct {
	db *gorm.DB
}

// NewGORMTransactionRepository creates a GORM-backed TransactionRepository.
func NewGORMTransactionRepository(db *gorm.DB) TransactionRepository {
	return &gormTransactionRepository{db: db}
}

func (r *gormTransactionRepository) Create(tx *domain.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *gormTransactionRepository) GetByID(id uuid.UUID) (*domain.Transaction, error) {
	var tx domain.Transaction
	if err := r.db.First(&tx, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *gormTransactionRepository) List(f Filter) ([]domain.Transaction, int64, error) {
	q := r.db.Model(&domain.Transaction{})

	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.Category != "" {
		q = q.Where("category = ?", f.Category)
	}
	if f.DateFrom != nil {
		q = q.Where("transaction_date >= ?", *f.DateFrom)
	}
	if f.DateTo != nil {
		// Extend to 23:59:59 so the entire end day is included.
		const endOfDay = 24*time.Hour - time.Second
		end := f.DateTo.Add(endOfDay)
		q = q.Where("transaction_date <= ?", end)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var txs []domain.Transaction
	if err := q.Order("transaction_date DESC").
		Limit(f.Limit).
		Offset(f.Offset).
		Find(&txs).Error; err != nil {
		return nil, 0, err
	}
	return txs, total, nil
}

func (r *gormTransactionRepository) Update(tx *domain.Transaction) error {
	return r.db.Save(tx).Error
}

func (r *gormTransactionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Transaction{}, "id = ?", id).Error
}

func (r *gormTransactionRepository) Summary() (*domain.Summary, error) {
	type result struct {
		Type  string
		Total float64
	}

	var rows []result
	if err := r.db.Model(&domain.Transaction{}).
		Select("type, SUM(amount) as total").
		Group("type").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	s := &domain.Summary{}
	for _, row := range rows {
		switch domain.TransactionType(row.Type) {
		case domain.TypeIncome:
			s.TotalIncome = row.Total
		case domain.TypeExpense:
			s.TotalExpense = row.Total
		}
	}
	s.Balance = s.TotalIncome - s.TotalExpense
	return s, nil
}
