package servicekit

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Endpoints collects all of the endpoints that compose a engine service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	RegisterEngineEndpoint      endpoint.Endpoint
	GetRegisteredEngineEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a enginesvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		RegisterEngineEndpoint:      MakeRegisterEngineEndpoint(s),
		GetRegisteredEngineEndpoint: MakeGetRegisteredEngineEndpoint(s),
	}
}

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a enginesvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path. That's fine: we simply need to provide specific encoders for
	// each endpoint.

	return Endpoints{
		RegisterEngineEndpoint:      httptransport.NewClient("POST", tgt, encodeRegisterEngineRequest, decodeRegisterEngineResponse, options...).Endpoint(),
		GetRegisteredEngineEndpoint: httptransport.NewClient("GET", tgt, encodeGetRegisteredEngineRequest, decodeGetRegisteredEngineResponse, options...).Endpoint(),
	}, nil
}

// RegisterEngine implements Service. Primarily useful in a client.
func (e Endpoints) RegisterEngine(ctx context.Context, eg Engine) error {
	request := registerEngineRequest{Engine: eg}
	response, err := e.RegisterEngineEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(registerEngineResponse)
	return resp.Err
}

// GetRegisteredEngine implements Service. Primarily useful in a client.
func (e Endpoints) GetRegisteredEngine(ctx context.Context, id string) (Engine, error) {
	request := getRegisteredEngineRequest{ID: id}
	response, err := e.GetRegisteredEngineEndpoint(ctx, request)
	if err != nil {
		return Engine{}, err
	}
	resp := response.(getRegisteredEngineResponse)
	return resp.Engine, resp.Err
}

// MakeRegisterEngineEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeRegisterEngineEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerEngineRequest)
		e := s.RegisterEngine(ctx, req.Engine)
		return registerEngineResponse{Err: e}, nil
	}
}

// MakeGetRegisteredEngineEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeGetRegisteredEngineEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRegisteredEngineRequest)
		eg, e := s.GetRegisteredEngine(ctx, req.ID)
		return getRegisteredEngineResponse{Engine: eg, Err: e}, nil
	}
}

// We have two options to return errors from the business logic.
//
// We could return the error via the endpoint itself. That makes certain things
// a little bit easier, like providing non-200 HTTP responses to the client. But
// Go kit assumes that endpoint errors are (or may be treated as)
// transport-domain errors. For example, an endpoint error will count against a
// circuit breaker error count.
//
// Therefore, it's often better to return service (business logic) errors in the
// response object. This means we have to do a bit more work in the HTTP
// response encoder to detect e.g. a not-found error and provide a proper HTTP
// status code. That work is done with the errorer interface, in transport.go.
// Response types that may contain business-logic errors implement that
// interface.

type registerEngineRequest struct {
	Engine Engine
}

type registerEngineResponse struct {
	Err error `json:"err,omitempty"`
}

func (r registerEngineResponse) error() error { return r.Err }

type getRegisteredEngineRequest struct {
	ID string
}

type getRegisteredEngineResponse struct {
	Engine Engine `json:"engine,omitempty"`
	Err    error  `json:"err,omitempty"`
}

func (r getRegisteredEngineResponse) error() error { return r.Err }