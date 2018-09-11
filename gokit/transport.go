package gokit

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	gokitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

// User represents a generic error
//
// swagger:model errorResponse
type ErrorResponse struct {
	// error message
	//
	// required: true
	Error string `json:"error"`
}

// EncodeError response back to the client
func EncodeError(c context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(&ErrorResponse{
		Error: err.Error(),
	})

}

/*
DecodeString will return the the body of the http request as a string
*/
func DecodeString(ctx context.Context, r *http.Request) (interface{}, error) {
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return string(s), nil
}

/*
DecodeResponse will check the response for an error, and if there is, it will set the body to the error message
*/
func DecodeResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != 200 {
		s, _ := ioutil.ReadAll(r.Body)
		return nil, errors.New(string(s))
	}
	return nil, nil
}

/*
EncodeResponse Convert the response into the response body
*/
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

/*
EncodeRequest converts the input request into a json string and adds it to the request body
*/
func EncodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(string)
	r.Body = nopCloser{strings.NewReader(req)}

	return nil
}

/*
EncodeToBase64 takes the byte array request and converts it to base64 before adding it to the request body
*/
func EncodeToBase64(_ context.Context, r *http.Request, request interface{}) error {
	req := request.([]byte)

	b := &bytes.Buffer{}
	e := base64.NewEncoder(base64.StdEncoding, b)
	e.Write(req)
	e.Close()

	r.Body = ioutil.NopCloser(b)
	return nil
}

func DecodeFromBase64(ctx context.Context, r *http.Request) (interface{}, error) {
	base64msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(base64msg))

}

/*
EncodeErrorRequest will extract the string message from the request error and add it to the body
*/
func EncodeErrorRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(error)
	r.Body = ioutil.NopCloser(bytes.NewReader([]byte(req.Error())))
	return nil
}

//EncodeJsonRequest will attempt to conver the input request
func EncodeJsonRequest(_ context.Context, r *http.Request, request interface{}) error {
	p := new(bytes.Buffer)
	err := json.NewEncoder(p).Encode(request)
	if err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(p)
	return nil
}

type errorer interface {
	error() error
}

func preProcess(s machineLearning.Service) kithttp.RequestFunc {
	return func(i context.Context, request *http.Request) context.Context {
		if s != nil {
			i = s.StoreHTTPRequest(i, request)
		}

		vars := mux.Vars(request)
		i = context.WithValue(i, "CHAT-ID", vars["chatid"])
		return i
	}
}

func GetChatId(ctx context.Context) uint32 {
	s := ctx.Value("CHAT-ID").(string)
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint32(i)
}

/*
GetServerOpts creates a default server option with an error logger, error encoder and a http request logger.
*/
func GetServerOpts(logger kitlog.Logger, service machineLearning.Service) []kithttp.ServerOption {
	return []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(EncodeError),
		kithttp.ServerBefore(gokitjwt.HTTPToContext(), preProcess(service)),
	}
}

func GetJWTKeys() jwt.Keyfunc {
	key := []byte(os.Getenv("JWT_KEY"))
	return func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}

}
