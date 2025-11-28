package domain

var ErrDrivingRecordNotFound DrivingRecordNotFoundError

type DrivingRecordNotFoundError struct{}

var _ error = DrivingRecordNotFoundError{}

func (DrivingRecordNotFoundError) Error() string { return "DrivingRecord not found" }

func (DrivingRecordNotFoundError) Is(other error) bool {
	_, ok := other.(DrivingRecordNotFoundError)
	return ok
}
