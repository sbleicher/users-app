package handler

import "github.com/labstack/echo/v4"

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func respError(c echo.Context, code int, message, details string) error {
	h := HttpError{
		Code:    code,
		Message: message,
		Details: details,
	}

	return c.JSON(code, h)
}

type HttpSuccess struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func respSuccess(c echo.Context, code int, message string, data ...interface{}) error {
	h := HttpSuccess{
		Code:    code,
		Message: message,
	}

	if len(data) > 0 {
		h.Data = data[0]
	}

	return c.JSON(code, h)
}
