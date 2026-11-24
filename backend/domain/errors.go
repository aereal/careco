package domain

var ErrDrivingRecordNotFound DrivingRecordNotFoundError

type DrivingRecordNotFoundError struct{}

var _ error = DrivingRecordNotFoundError{}

func (DrivingRecordNotFoundError) Error() string { return "DrivingRecord not found" }

func (DrivingRecordNotFoundError) Is(other error) bool {
	_, ok := other.(DrivingRecordNotFoundError)
	return ok
}

var ErrTireNotFound TireNotFoundError

type TireNotFoundError struct{}

var _ error = TireNotFoundError{}

func (TireNotFoundError) Error() string { return "Tire not found" }

func (TireNotFoundError) Is(other error) bool {
	_, ok := other.(TireNotFoundError)
	return ok
}

var ErrTireRetired TireRetiredError

type TireRetiredError struct{}

var _ error = TireRetiredError{}

func (TireRetiredError) Error() string { return "Tire is retired" }

func (TireRetiredError) Is(other error) bool {
	_, ok := other.(TireRetiredError)
	return ok
}
