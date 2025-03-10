package customerror

import "fmt"

type ErrorableInterface interface {
	Code() int
	Message() string
	error
}


type BaseError struct {
	ErrMessage string
	ErrCode int
}


func (e BaseError) Error() string{
	return fmt.Sprintf("Error: %s (Code: %d)", e.ErrMessage, e.ErrCode)
}

func (e BaseError) Code() int{
	return e.ErrCode
}

func (e BaseError) Message() string{
	return e.ErrMessage
}


type ServerError struct {
	ErrMessage string
	ErrCode int
}


func (e *ServerError) Error() string{
	return fmt.Sprintf("Error: %s (Code: %d)", e.ErrMessage, e.ErrCode)
}

func (e *ServerError) Code() int{
	return e.ErrCode
}

func (e *ServerError) Message() string{
	return e.ErrMessage
}


type ValidationError struct {
	ErrMessage string
	ErrCode int
	ErrBag map[string][]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Error: %s (Code: %d)", e.ErrMessage, e.ErrCode)
}

func (e *ValidationError) Code() int{
	return e.ErrCode
}

func (e *ValidationError) Message() string{
	return e.ErrMessage
}


func (e *ValidationError) ErrorBag() map[string][]string{
	return e.ErrBag
}
