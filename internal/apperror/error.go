package apperror

import "encoding/json"

var (
	ErrorNotFound = NewAppError(nil, "not found", "", "US-000003")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"massage,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, massage, DeveloperMessage, Code string) *AppError {
	return &AppError{
		Err:              err,
		Message:          massage,
		DeveloperMessage: DeveloperMessage,
		Code:             Code,
	}
}

func SystemError(err error) *AppError {
	return NewAppError(err, "SystemError", err.Error(), "US-000000")
}
