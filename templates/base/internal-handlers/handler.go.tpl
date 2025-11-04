// Package handlers provide way to work with net/http package
package handlers

import (
	"net/http"
	"reflect"
	"strings"

	"{{.ProjectName}}/cmd/config"
)

var H Handler

// any metadata
type Metadata struct{}

type Handler struct {
	app *config.Application
}

func NewHandler(app *config.Application) {
	H = Handler{
		app: app,
	}
}

func (h *Handler) GetFilterableFields(r *http.Request) map[string][]string {
	query := r.URL.Query()
	meta := Metadata{}

	//reflect metadata struct
	ref := reflect.ValueOf(meta)
	refType := ref.Type()
	//loop through metadata fields
	for i := 0; i < refType.NumField(); i++ {
		if splited := strings.Split(refType.Field(i).Tag.Get("json"), ",")[0]; splited != "" {
			query.Del(splited)
		}
	}

	query.Del("locale")

	return query
}

