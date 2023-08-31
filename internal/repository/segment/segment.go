package segment

import (
	"fmt"

	"gorm.io/gorm"
)

type Segment struct {
	Slug string `gorm:"column:slug;primaryKey"`
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*Repository, error) {
	if !db.Migrator().HasTable(&Segment{}) {
		err := db.Migrator().AutoMigrate(&Segment{})
		if err != nil {
			fmt.Println(err)
		}
	}
	return &Repository{
		db: db,
	}, nil
}

func (s *Repository) CreateSegment(slug string) error {
	segment := Segment{
		Slug: slug,
	}
	tx := s.db.Create(&segment)
	if tx.Error != nil {
		fmt.Printf("failed to add %s\n", slug)
		return tx.Error
	}
	fmt.Printf("%s successfully added\n", slug)
	return nil
}

func (s *Repository) DeleteSegment(slug string) error {
	segment := Segment{
		Slug: slug,
	}
	tx := s.db.Delete(segment)
	if tx.Error != nil {
		fmt.Printf("failed to delete %s\n", slug)
		return tx.Error
	}
	fmt.Printf("%s successfully deleted\n", slug)
	return nil
}
