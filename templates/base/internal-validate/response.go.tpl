package validate

import (
	"errors"
	"net/http"

	"{{.ProjectName}}/cmd/config"

	"github.com/go-playground/validator/v10"
)

var app *config.Application

func NewValidate(a *config.Application) {
	app = a
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorLog.Printf("internal server error %s, %s, %s", r.Method, r.URL, err.Error())
	WriteJSONError(w, http.StatusInternalServerError, "internal server error", nil)
}

func UnauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorLog.Printf("unauthorized %s, %s, %s", r.Method, r.URL, err.Error())
	WriteJSONError(w, http.StatusUnauthorized, "unauthorized", nil)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorLog.Printf("bad request error %s, %s, %s", r.Method, r.URL, err.Error())
	WriteJSONError(w, http.StatusBadRequest, err.Error(), nil)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorLog.Printf("not found error %s, %s, %s", r.Method, r.URL, err.Error())
	WriteJSONError(w, http.StatusNotFound, "not found", nil)
}

func ErrConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorLog.Printf("conflict error %s, %s, %s", r.Method, r.URL, err.Error())
	WriteJSONError(w, http.StatusConflict, err.Error(), nil)
}

func JsonValidationErrorResponse(w http.ResponseWriter, err error) error {
	var errs validator.ValidationErrors
	errors.As(err, &errs)
	transformedErrors := RangeErrors(errs)
	return WriteJSONError(w, http.StatusUnprocessableEntity, "unprocessable entity", transformedErrors)
}

func SendResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	if err := JsonResponse(w, status, data); err != nil {
		InternalServerError(w, r, err)
		return
	}
}

