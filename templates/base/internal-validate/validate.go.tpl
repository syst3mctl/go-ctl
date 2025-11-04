// Package validate is used to validate struct fields, read and write json
package validate

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

type ErrorEnvelope struct {
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}
type Envelope struct {
	Data ErrorEnvelope `json:"data"`
}

func WriteJSONError(w http.ResponseWriter, status int, message string, errors any) error {

	return WriteJSON(w, status, Envelope{
		Data: ErrorEnvelope{
			Message: message,
			Errors:  errors,
		},
	})
}

func RangeErrors(err validator.ValidationErrors) *map[string][]string {
	errs := map[string][]string{}

	for _, e := range err {
		key := strings.ToLower(e.Field())
		errs[key] = append(errs[key], e.Error())
	}

	return &errs
}

func JsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return WriteJSON(w, status, envelope{Data: data})
}

