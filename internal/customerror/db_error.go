package customerror

import "fmt"

type DBReadError struct {
	ErrMessage string
	ErrCode int
}


func (e DBReadError) Error() string{
	return fmt.Sprintf("Error: %s (Code: %d)", e.ErrMessage, e.ErrCode)
}

func (e DBReadError) Code() int{
	return e.ErrCode
}

func (e DBReadError) Message() string{
	return e.ErrMessage
}
