package ecpay

import "fmt"

// ErrorCode 錯誤代碼
type ErrorCode string

const (
	ErrCodeValidation ErrorCode = "VALIDATION_ERROR"
	ErrCodeRequest    ErrorCode = "REQUEST_ERROR"
	ErrCodeNetwork    ErrorCode = "NETWORK_ERROR"
	ErrCodeResponse   ErrorCode = "RESPONSE_ERROR"
	ErrCodeParse      ErrorCode = "PARSE_ERROR"
	ErrCodeAPI        ErrorCode = "API_ERROR"
	ErrCodeCrypto     ErrorCode = "CRYPTO_ERROR"
)

// Error 自定義錯誤
type Error struct {
	Code    ErrorCode
	Message string
}

// NewError 建立新的錯誤
func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error 實作 error 介面
func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// IsError 檢查是否為特定錯誤類型
func IsError(err error, code ErrorCode) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}