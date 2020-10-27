// Package client provides a enginesvc client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package client

import (
	"io"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/isgo-golgo13/go-gokit-gorilla-restsvc/engine_svcpkg"
)

// New returns a service that's load-balanced over instances of enginesvc found
// in the provided Consul server. The mechanism of looking up emginesvc
// instances in Consul is hard-coded into the client.
func New(consulAddr string, logger log.Logger) (engine_svcpkg.Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	// As the implementer of enginesvc, we declare and enforce these
	// parameters for all of the enginesvc consumers.
	var (
		consulService = "enginesvc"
		consulTags    = []string{"prod"}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiclient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		endpoints engine_svcpkg.Endpoints
	)
	{
		factory := factoryFor(engine_svcpkg.MakeRegisterEngineEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.RegisterEngineEndpoint = retry
	}
	{
		factory := factoryFor(engine_svcpkg.MakeGetRegisteredEngineEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetRegisteredEngineEndpoint = retry
	}
	
	return endpoints, nil
}

func factoryFor(makeEndpoint func(engine_svcpkg.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := engine_svcpkg.MakeClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}