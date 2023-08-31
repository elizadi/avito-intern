package operation

import (
	"avito/internal/domain"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Operation struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	UserID      uint64    `gorm:"column:user_id"`
	SegmentSlug string    `gorm:"column:segment_slug"`
	Operation   string    `gorm:"column:operation"`
	Data        time.Time `gorm:"column:data"`
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*Repository, error) {
	if !db.Migrator().HasTable(&Operation{}) {
		err := db.Migrator().AutoMigrate(&Operation{})
		if err != nil {
			fmt.Println(err)
		}
	}
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) CreateOperationBatch(slugsAdd, slugsDelete []string, id uint64) error {
	var operations = createBatch(slugsAdd, id, "Add")
	if len(operations) != 0 {
		tx := r.db.Create(&operations)
		if tx.Error != nil {
			return tx.Error
		}
	}
	
	operations = createBatch(slugsDelete, id, "Delete")
	if len(operations) != 0 {
		tx := r.db.Create(&operations)
		if tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}

func (r *Repository) GetOperations(year, month int) ([]domain.Operation, error) {
	var operations []Operation
	err := r.db.Debug().
		Model(&Operation{}).
		Where("EXTRACT('Year' FROM \"data\") = ? AND EXTRACT('Month' FROM \"data\") = ?", year, month).
		Find(&operations).
		Error
	if err != nil {
		return nil, err
	}
	return toDomainOperations(operations), nil
}
func toDomainOperations(operations []Operation) []domain.Operation {
	res := make([]domain.Operation, 0, len(operations))
	for _, operation := range operations {
		res = append(res, domain.Operation{
			ID:          operation.ID,
			UserID:      operation.UserID,
			SegmentSlug: operation.SegmentSlug,
			Operation:   operation.Operation,
			Data:        operation.Data,
		})
	}
	return res
}
func createBatch(slugs []string, id uint64, operation string) []Operation {
	res := make([]Operation, 0, len(slugs))
	now := time.Now()
	for _, slug := range slugs {
		res = append(res, Operation{
			UserID:      id,
			SegmentSlug: slug,
			Operation:   operation,
			Data:        now,
		})
	}
	return res
}
