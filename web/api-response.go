package web

import "fmt"

type DataDto struct {
	Type    string      `json:"@type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func dataResponse(responseType string, message string, data interface{}) DataDto {
	return DataDto{
		Type:    responseType,
		Message: message,
		Data:    data,
	}
}

type ErrorDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ErrorDto) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func errorDto(code, message string) ErrorDto {
	return ErrorDto{
		Code:    code,
		Message: message,
	}
}
