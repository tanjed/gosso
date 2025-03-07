package customerror

import "fmt"

type UserNotFoundError struct {
	ErrMessage string
	ErrCode int
}

func (e UserNotFoundError) Error() string{
	return fmt.Sprintf("Error: %s (Code: %d)", e.ErrMessage, e.ErrCode)
}

func (e UserNotFoundError) Code() int{
	return e.ErrCode
}

func (e UserNotFoundError) Message() string{
	return e.ErrMessage
}



