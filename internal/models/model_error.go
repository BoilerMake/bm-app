package models

type ModelError struct {
	Message string
	Type    int
}

func (e ModelError) Error() string {
	return e.Message
}

func (e ModelError) GetType() int {
	return e.Type
}
