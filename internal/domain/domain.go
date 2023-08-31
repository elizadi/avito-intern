package domain

type UseCase interface {
	CreateSegment(slug string) error
	DeleteSegment(slug string) error
	AddUserToSegment(slugsAdd []string, slugsDelete []string, id uint64) error
	GetActiveUserSegments(id uint64) ([]Segment, error)
	GetOperations(year, month int) ([]Operation, error)
}

type OperationRepository interface {
	CreateOperationBatch(slugsAdd, slugsDelete []string, id uint64) error
	GetOperations(year, month int) ([]Operation, error)
}

type SegmentRepository interface {
	CreateSegment(slug string) error
	DeleteSegment(slug string) error
}

type UserSegmentRepository interface {
	AddUserToSegment(slugsAdd []string, slugsDelete []string, id uint64) ([]string, []string, error)
	GetActiveUserSegments(id uint64) ([]Segment, error)
}
