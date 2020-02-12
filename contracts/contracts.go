package contracts

import (
	"fmt"
)

const (
	notExistsMessage   = "Value with key='%s' doesn't exist."
	keyToLongMessage   = "Provided key size (%d bytes) exceeds maximum allowed key size (%d bytes)."
	valueToLongMessage = "Provided value size (%d bytes) exceeds maximum allowed value size (%d bytes)."
)

type Error struct {
	Error string `json:"error"`
}

func NotExistsError(key string) *Error {
	return &Error{Error: fmt.Sprintf(notExistsMessage, key)}
}

func KeyToLongError(actualKey, maxKey int) *Error {
	return &Error{Error: fmt.Sprintf(keyToLongMessage, actualKey, maxKey)}
}

func ValueToLongError(actualValue, maxValue int) *Error {
	return &Error{Error: fmt.Sprintf(valueToLongMessage, actualValue, maxValue)}
}
