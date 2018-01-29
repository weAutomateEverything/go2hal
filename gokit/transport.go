package gokit

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	http2 "github.com/go-kit/kit/transport/http"
	"log"
	"fmt"
	"strings"
	"bytes"
)

// encode errors from business-logic
func EncodeError(c context.Context, err error, w http.ResponseWriter) {
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

func LogRequest() http2.RequestFunc{
	return func(i context.Context, request *http.Request) context.Context {
		log.Print(formatRequest(request))
		return i
	}
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
	r.ParseForm()
	request = append(request, "\n")
	request = append(request, r.Form.Encode())
	}

	body,_ := ioutil.ReadAll(r.Body)
	request = append(request, fmt.Sprintf("Body: %v",string(body)))

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Return the request as a string
	return strings.Join(request, "\n")
}
