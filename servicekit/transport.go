package servicekit

// The enginesvcpkg is just over HTTP, so we just have a single transport.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrErrorInRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrErrorInRouting = errors.New("Error in the routing")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a enginesvc server.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST    /engines/                          registers another engine
	// GET     /engines/:id                       retrieves the given engine by id

	r.Methods("POST").Path("/engines/").Handler(httptransport.NewServer(
		e.RegisterEngineEndpoint,
		decodeRegisterEngineRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/engines/{id}").Handler(httptransport.NewServer(
		e.GetRegisteredEngineEndpoint,
		decodeGetRegisteredEngineRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeRegisterEngineRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req registerEngineRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Engine); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetRegisteredEngineRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrErrorInRouting
	}
	return getRegisteredEngineRequest{ID: id}, nil
}

func encodeRegisterEngineRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/engines/")
	req.URL.Path = "/engines/"
	return encodeRequest(ctx, req, request)
}

func encodeGetRegisteredEngineRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/engines/{id}")
	r := request.(getRegisteredEngineRequest)
	engineID := url.QueryEscape(r.ID)
	req.URL.Path = "/engines/" + engineID
	return encodeRequest(ctx, req, request)
}

func decodeRegisterEngineResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response registerEngineResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetRegisteredEngineResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getRegisteredEngineResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// enginesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrEngineNotExistInRegistry:
		return http.StatusNotFound
	case ErrEnginePreExistInRegistry, ErrInconsistentEngineIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}