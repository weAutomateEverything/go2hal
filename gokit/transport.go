package gokit

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// encode errors from business-logic
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})

}

func DecodeString(_ context.Context, r *http.Request) (interface{}, error) {
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return string(s),nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if response == nil {
		return nil
	}
	if e, ok := response.(errorer); ok && e.error() != nil {
		EncodeError(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}
