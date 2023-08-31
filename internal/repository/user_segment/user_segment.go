package user_segment

import (
	"avito/internal/domain"
	"avito/internal/repository/segment"
	"avito/internal/repository/user"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserSegment struct {
	UserID  uint64          `gorm:"column:user_id;primaryKey;"`
	User    user.User       `gorm:"foreignKey:user_id"`
	Slug    string          `gorm:"column:slug;primaryKey;"`
	Segment segment.Segment `gorm:"foreignKey:slug"`
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*Repository, error) {
	if !db.Migrator().HasTable(&UserSegment{}) {
		err := db.Migrator().AutoMigrate(&UserSegment{})
		if err != nil {
			fmt.Println(err)
		}
	}
	return &Repository{
		db: db,
	}, nil
}

func (s *Repository) AddUserToSegment(slugsAdd []string, slugsDelete []string, id uint64) ([]string, []string, error) {
	segmentsAdd := createBatch(slugsAdd, id)
	if len(segmentsAdd) != 0 {
		tx := s.db.
			Clauses(clause.OnConflict{DoNothing: true}).
			Clauses(clause.Returning{}).
			Model(&UserSegment{}).
			Create(&segmentsAdd)
		if tx.Error != nil {
			return nil, nil, tx.Error
		}
	}

	segmentsDelete := createBatch(slugsDelete, id)
	if len(segmentsDelete) != 0 {
		tx := s.db.
			Clauses(clause.Returning{}).
			Model(&UserSegment{}).
			Delete(&segmentsDelete)
		if tx.Error != nil {
			return nil, nil, tx.Error
		}
	}
	return slugs(segmentsAdd), slugs(segmentsDelete), nil
}

func (s *Repository) GetActiveUserSegments(id uint64) ([]domain.Segment, error) {
	var res []*UserSegment

	tx := s.db.Model(&UserSegment{}).Preload(clause.Associations).Where("user_id = ?", id).Find(&res)
	err := tx.Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(res)

	return todomainSegments(res), nil
}

func createBatch(slugs []string, id uint64) []*UserSegment {
	res := make([]*UserSegment, 0, len(slugs))
	for _, slug := range slugs {
		res = append(res, &UserSegment{
			UserID: id,
			Slug:   slug,
		})
	}
	return res
}

func todomainSegments(segments []*UserSegment) []domain.Segment {
	res := make([]domain.Segment, 0, len(segments))
	for _, segment := range segments {
		if segment == nil {
			continue
		}
		res = append(res, domain.Segment{
			Slug: segment.Slug,
		})
	}
	return res
}

func slugs(segments []*UserSegment) []string {
	res := make([]string, 0, len(segments))
	for _, segment := range segments {
		if segment == nil || segment.Slug == "" || segment.UserID == 0 {
			continue
		}
		res = append(res, segment.Slug)
	}
	return res
}

// func createBatch(slugs []string, id uint64, operation string) []Operation {
// 	res := make([]Operation, 0, len(slugs))
// 	now := time.Now()
// 	for _, slug := range slugs {
// 		res = append(res, Operation{
// 			UserID:      id,
// 			SegmentSlug: slug,
// 			Operation:   operation,
// 			Data:        now,
// 		})
// 	}
// 	return res
// }
