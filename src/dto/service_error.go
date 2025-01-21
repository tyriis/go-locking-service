package dto

type ServiceError struct {
	Service string
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	return "ServiceError:" + e.Service + ":" + e.Message + "\n" + e.Err.Error()
}
