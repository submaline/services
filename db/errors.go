package db

import "fmt"

func NewInvalidArgumentError(paramName string) error {
	return fmt.Errorf("invalid argument: %v", paramName)
}
