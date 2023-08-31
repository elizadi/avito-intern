package app

import (
	"avito/internal/domain"
	"avito/internal/repository/operation"
	"avito/internal/repository/segment"
	"avito/internal/repository/user"
	"avito/internal/repository/user_segment"
	"avito/internal/usecase"
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Execute() (domain.UseCase, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, errors.New("empty parameter")
	}
	dsn := fmt.Sprintf("host=%s user=postgres password=123456789Lis port=5432 sslmode=disable", url)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	repoSegment, err := segment.New(db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_, err = user.New(db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	repoUserSegment, err := user_segment.New(db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	repoOperation, err := operation.New(db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	uc, err := usecase.New(repoSegment, repoOperation, repoUserSegment)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return uc, nil
}
