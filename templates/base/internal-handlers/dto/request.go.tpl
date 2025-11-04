package dto

import (
	"net/http"
	"reflect"

	"{{.ProjectName}}/internal/validate"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func ParseAndValidate(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	err := validate.ReadJSON(w, r, dest)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Slice {
		slice := v.Elem()
		for i := 0; i < slice.Len(); i++ {
			item := slice.Index(i).Interface()
			if err = Validate.Struct(item); err != nil {
				return err
			}
		}
		return err
	}

	return Validate.Struct(dest)
}

