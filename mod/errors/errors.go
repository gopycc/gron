package errors

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	Success ErrorCode = 200
	Unkown            = -1
)

func New(text string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(text, args))
}

type ArgumentNilOrEmpty struct {
	Name string
}

func (err ArgumentNilOrEmpty) Error() string {
	return fmt.Sprintf("%s is nil", err.Name)
}
