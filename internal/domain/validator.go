package domain

type Validator interface {
	Validate(data interface{}) error
}

type ValidationRepository interface {
	ValidateSchema(data interface{}) error
}
