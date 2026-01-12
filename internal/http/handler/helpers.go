package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"Beanefits/internal/api"
)

type Helpers struct {
	JSON func(w http.ResponseWriter, status int, v any)
}

func DefaultHelpers() Helpers {
	return Helpers{
		JSON: func(w http.ResponseWriter, status int, v any) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(v)
		},
	}
}

func WriteProblem(w http.ResponseWriter, status int, title string, detail *string, code *string, instance *string) {
	p := api.Problem{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Code:     code,
		Instance: instance,
	}

	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(p)
}

func instanceFromRequest(r *http.Request) *string {
	if r == nil || r.URL == nil || r.URL.Path == "" {
		return nil
	}
	p := r.URL.Path
	return &p
}

func ptr[T any](v T) *T { return &v }

func DecodeJSON[T any](r *http.Request, dst *T) error {
	if r == nil || r.Body == nil {
		return errors.New("empty body")
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	if err := dec.Decode(new(any)); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	return errors.New("unexpected extra json content")
}
