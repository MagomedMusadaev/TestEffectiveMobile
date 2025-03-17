package entities

// ErrorResponse описывает структуру ошибки
// @Description Структура для представления ошибки, которая включает код ошибки и сообщение.
// @Example { "code": 500, "message": "Internal Server Error" }
type ErrorResponse struct {
	// Код ошибки
	// @example 500
	Code int `json:"code"`

	// Сообщение ошибки
	// @example "Internal Server Error"
	Message string `json:"message"`
}
