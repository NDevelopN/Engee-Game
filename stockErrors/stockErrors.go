package stockErrors

import (
	"fmt"
)

type EmptyValueError struct {
	Field string
}

func (e EmptyValueError) Error() string {
	return fmt.Sprintf("no value provided for field %s", e.Field)
}

type InvalidValueError[T any] struct {
	Field string
	Value T
}

func (e InvalidValueError[T]) Error() string {
	return fmt.Sprintf("invalid value '%v' for field %s", e.Value, e.Field)
}

type MatchNotFoundError[T any] struct {
	Space string
	Field string
	Value T
}

func (e *MatchNotFoundError[T]) Error() string {
	if e.Space == "" {
		return fmt.Sprintf("no match found for provided %s: %v", e.Field, e.Value)
	}

	return fmt.Sprintf("no match found in %s for provided %s: %v", e.Space, e.Field, e.Value)
}

type MatchFoundError[T any] struct {
	Space string
	Field string
	Value T
}

func (e *MatchFoundError[T]) Error() string {
	if e.Space == "" {
		return fmt.Sprintf("existing match found for provided %s: %v", e.Field, e.Value)

	}

	return fmt.Sprintf("%s already exists in %s for value: %v", e.Field, e.Space, e.Value)
}

type EmptySetError struct {
	Set string
}

func (e *EmptySetError) Error() string {
	return fmt.Sprintf("set %s is empty", e.Set)
}

type InvalidSetLengthError struct {
	Set      string
	Expected int
	Actual   int
}

func (e *InvalidSetLengthError) Error() string {
	return fmt.Sprintf("set %s is incorrect length '%d', expected '%d'", e.Set, e.Actual, e.Expected)
}

type SingleSendError struct {
	Message string
	Target  string
	Err     error
}

func (e *SingleSendError) Error() string {
	return fmt.Sprintf("could not send %q message to %s: %v", e.Message, e.Target, e.Err)
}

type BroadcastError struct {
	Message string
	Err     error
}

func (e *BroadcastError) Error() string {
	return fmt.Sprintf("could not broadcast %q message: %v", e.Message, e.Err)
}
