package commons

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApiResponseBody[T any] struct {
	ApiResponseBase
	Data T `json:"data"`
}

type ApiResponseBase struct {
	Status  string `json:"status" example:"success"`
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"success"`
}

type ApiResponsePagination struct {
	Count       int             `json:"count"`
	CurrentPage int             `json:"current_page"`
	Link        ApiResponseLink `json:"link"`
	PerPage     int             `json:"per_page"`
	Total       int             `json:"total"`
	TotalPages  int             `json:"total_pages"`
}

type ApiResponseLink struct {
	Next     *string `json:"next,omitempty"`
	Previous *string `json:"previous,omitempty"`
}

type ApiResponseMeta struct {
	Pagination ApiResponsePagination `json:"pagination"`
}

type ApiResponseData[T any] struct {
	ApiResponseBase
	Data T               `json:"data"`
	Meta ApiResponseMeta `json:"meta"`
}

type ApiResponseMetaData[T any] struct {
	ApiResponseBase
	Data T              `json:"data"`
	Meta map[string]any `json:"meta"`
}

type ApiResponseTotalData[T any] struct {
	ApiResponseBase
	Data T                     `json:"data"`
	Meta ApiResponseAmountMeta `json:"meta"`
}

type ApiResponseAmountMeta struct {
	Pagination  ApiResponsePagination `json:"pagination"`
	TotalAmount int                   `json:"total_amount"`
}

func ResponseSuccess[T any](message string, data T) ApiResponseBody[T] {
	return ApiResponseBody[T]{
		ApiResponseBase: ApiResponseBase{
			Status:  "success",
			Code:    http.StatusOK,
			Message: message,
		},
		Data: data,
	}
}

func ResponseSuccessWithoutMessage[T any](data T) ApiResponseBody[T] {
	return ApiResponseBody[T]{
		ApiResponseBase: ApiResponseBase{
			Status:  "success",
			Code:    http.StatusOK,
			Message: "OK",
		},
		Data: data,
	}
}

func ResponseSuccessWithoutData(message string) ApiResponseBase {
	return ApiResponseBase{
		Status:  "success",
		Code:    http.StatusOK,
		Message: message,
	}
}

func ResponseFailedServerError(message string) ApiResponseBase {
	return ApiResponseBase{
		Status:  "failed",
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func ResponseFailed(message string) ApiResponseBase {
	return ApiResponseBase{
		Status:  "failed",
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func ResponseFailedCode(code int, message string) ApiResponseBase {
	return ApiResponseBase{
		Status:  "failed",
		Code:    code,
		Message: message,
	}
}

func HttpResponseFailed(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusBadRequest, ResponseFailed(message))
}

func HttpResponseSuccess[T any](ctx echo.Context, message T) error {
	return ctx.JSON(http.StatusOK, ResponseSuccess("success", message))
}

func HttpResponseCreated(ctx echo.Context, message interface{}) error {
	return ctx.JSON(http.StatusCreated, ResponseSuccess("success", message))
}

func MapErrorResponse(err error) ApiResponseBody[any] {
	// duplErr := commons.ErrDuplicate{}
	// if errors.As(err, &duplErr) {
	// 	msg := fmt.Sprintf("%s already exist", duplErr.Field)
	// 	return ResponseFailedWithData("duplicated", map[string]interface{}{
	// 		"field":   duplErr.Field,
	// 		"message": msg,
	// 	})
	// }
	return ApiResponseBody[any]{
		ApiResponseBase: ResponseFailed(err.Error()),
		Data:            nil,
	}

}

func ResponseFailedWithData(message string, data interface{}) ApiResponseBody[any] {
	return ApiResponseBody[any]{
		ApiResponseBase: ApiResponseBase{
			Message: message,
			Status:  "failed",
		},
		Data: data,
	}
}

func ResponseFailedWithError[T error, K string](message string, data T) ApiResponseBody[K] {
	return ApiResponseBody[K]{
		ApiResponseBase: ApiResponseBase{
			Status:  "failed",
			Message: message,
		},
		Data: K(data.Error()),
	}
}

func ResponseMetaDataSuccess[T any](
	message string,
	meta map[string]any,
	data T,
) ApiResponseMetaData[T] {
	return ApiResponseMetaData[T]{
		ApiResponseBase: ApiResponseBase{
			Status:  "success",
			Code:    http.StatusOK,
			Message: message,
		},
		Meta: meta,
		Data: data,
	}
}

func ResponseDataSuccess[T any](
	message string,
	pagination ApiResponsePagination,
	data T,
) ApiResponseData[T] {
	return ApiResponseData[T]{
		ApiResponseBase: ApiResponseBase{
			Status:  "success",
			Code:    http.StatusOK,
			Message: message,
		},
		Meta: ApiResponseMeta{
			Pagination: pagination,
		},
		Data: data,
	}
}

func ResponseAmountDataSuccess[T any](
	message string,
	pagination ApiResponsePagination,
	data T,
	total int,
) ApiResponseTotalData[T] {
	return ApiResponseTotalData[T]{
		ApiResponseBase: ApiResponseBase{
			Status:  "success",
			Code:    http.StatusOK,
			Message: message,
		},
		Meta: ApiResponseAmountMeta{
			Pagination:  pagination,
			TotalAmount: total,
		},
		Data: data,
	}
}

func ResponseUnauthorized(message string) ApiResponseBase {
	return ApiResponseBase{
		Status:  "failed",
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func ResponseRequired(message string) ApiResponseBase {
	return ApiResponseBase{
		Status:  "failed",
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

type Auth struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
