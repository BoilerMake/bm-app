package models

import "fmt"

type ModelError struct {
	Message string
}

func (e *ModelError) Error() string {
	return fmt.Sprintf("model error: %v", e.Message)
}
