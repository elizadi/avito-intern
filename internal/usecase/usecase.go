package usecase

import (
	"avito/internal/domain"
)

type usecase struct {
	repoSegment     domain.SegmentRepository
	repoOperation   domain.OperationRepository
	repoUserSegment domain.UserSegmentRepository
}

func New(
	repoSegment domain.SegmentRepository,
	repoOperation domain.OperationRepository,
	repoUserSegment domain.UserSegmentRepository,
) (*usecase, error) {
	return &usecase{
		repoSegment:     repoSegment,
		repoOperation:   repoOperation,
		repoUserSegment: repoUserSegment,
	}, nil
}

func (u *usecase) CreateSegment(slug string) error {
	err := u.repoSegment.CreateSegment(slug)
	if err != nil {
		return err
	}
	return nil
}

func (u *usecase) DeleteSegment(slug string) error {
	err := u.repoSegment.DeleteSegment(slug)
	if err != nil {
		return err
	}
	return nil
}

func (u *usecase) AddUserToSegment(slugsAdd []string, slugsDelete []string, id uint64) error {
	if len(slugsAdd) == 0 && len(slugsDelete) == 0 {
		return domain.ErrEmptyParameter
	}
	added, deleted, err := u.repoUserSegment.AddUserToSegment(slugsAdd, slugsDelete, id)
	if err != nil {
		return err
	}
	err = u.repoOperation.CreateOperationBatch(added, deleted, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *usecase) GetActiveUserSegments(id uint64) ([]domain.Segment, error) {
	segments, err := u.repoUserSegment.GetActiveUserSegments(id)
	if err != nil {
		return nil, err
	}
	return segments, nil
}

func (u *usecase) GetOperations(year, month int) ([]domain.Operation, error) {
	operations, err := u.repoOperation.GetOperations(year, month)
	if err != nil {
		return nil, err
	}
	return operations, nil
}
