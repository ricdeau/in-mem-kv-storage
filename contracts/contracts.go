package contracts

import (
	"fmt"
)

const (
	notExistsMessage = "Value with key='%s' doesn't exist"
)

type Error struct {
	Error string `json:"error"`
}

func NotExistsError(key string) *Error {
	return &Error{Error: fmt.Sprintf(notExistsMessage, key)}
}
