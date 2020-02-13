package contracts

import (
	"fmt"
)

// Message templates.
const (
	notExistsMessage   = "Value with key='%s' doesn't exist."
	keyToLongMessage   = "Provided key size (%d bytes) exceeds maximum allowed key size (%d bytes)."
	valueToLongMessage = "Provided value size (%d bytes) exceeds maximum allowed value size (%d bytes)."
)

// Error represents error understandable to the user.
type Error struct {
	Error string `json:"error"`
}

// NotExistsError returns the error that occurred when storage does't contains given key.
// Key will be included in the error's message.
func NotExistsError(key string) *Error {
	return &Error{Error: fmt.Sprintf(notExistsMessage, key)}
}

// KeyToLongError returns the error that occurred when given key exceeds key limit.
// Actual and maximum key lengths in bytes will be included in the error's message.
func KeyToLongError(actualKey, maxKey int) *Error {
	return &Error{Error: fmt.Sprintf(keyToLongMessage, actualKey, maxKey)}
}

// ValueToLongError returns the error that occurred when given value exceeds value limit.
// Actual and maximum value lengths in bytes will be included in the error's message.
func ValueToLongError(actualValue, maxValue int) *Error {
	return &Error{Error: fmt.Sprintf(valueToLongMessage, actualValue, maxValue)}
}
