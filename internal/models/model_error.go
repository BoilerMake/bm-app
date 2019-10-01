package models

type ModelError struct {
	Message string
}

func (e ModelError) Error() string {
	return e.Message
}
