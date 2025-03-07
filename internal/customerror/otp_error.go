package customerror

import "fmt"

type OtpNotFoundError struct {
	ErrMessage string
	ErrCode int
}


func (e *OtpNotFoundError) Error() string{
	return fmt.Sprintf("Error: %s (Code: %d)", e.ErrMessage, e.ErrCode)
}

func (e *OtpNotFoundError) Code() int{
	return e.ErrCode
}

func (e *OtpNotFoundError) Message() string{
	return e.ErrMessage
}


type OtpLimitReachedError struct {
	ErrMessage string
	ErrCode int
}


func (e OtpLimitReachedError) Error() string{
	return fmt.Sprintf("Error: %s (Code: %d)", e.ErrMessage, e.ErrCode)
}

func (e OtpLimitReachedError) Code() int{
	return e.ErrCode
}

func (e OtpLimitReachedError) Message() string{
	return e.ErrMessage
}




type OtpAlreadySentError struct {
	ErrMessage string
	ErrCode int
}


func (e OtpAlreadySentError) Error() string{
	return fmt.Sprintf("Error: %s (Code: %d)", e.ErrMessage, e.ErrCode)
}

func (e OtpAlreadySentError) Code() int{
	return e.ErrCode
}

func (e OtpAlreadySentError) Message() string{
	return e.ErrMessage
}




type OtpMismatchError struct {
	ErrMessage string
	ErrCode int
}


func (e *OtpMismatchError) Error() string{
	return fmt.Sprintf("Error: %s (Code: %d)", e.ErrMessage, e.ErrCode)
}

func (e *OtpMismatchError) Code() int{
	return e.ErrCode
}

func (e *OtpMismatchError) Message() string{
	return e.ErrMessage
}
