package dto

type DAOError struct {
	Operation string
	Err       error
}

func (e *DAOError) Error() string {
	return e.Operation + ": " + e.Err.Error()
}
