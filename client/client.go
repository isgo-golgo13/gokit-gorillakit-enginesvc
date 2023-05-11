// Package client provides a enginesvc client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package client

import (
	"io"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"gokit-gorillakit-enginesvc/servicekit"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
)

// New returns a service that's load-balanced over instances of enginesvc found
// in the provided Consul server. The mechanism of looking up ' engine service''
// instances in Consul is hard-coded into the client.
func New(consulAddr string, logger log.Logger) (servicekit.Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	// As the implementer of 'gokit-enginesvc', we declare and enforce these
	// parameters for all of the 'service' consumers.
	var (
		consulService = "gokit-enginesvc"
		consulTags    = []string{"prod"}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiclient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		endpoints servicekit.Endpoints
	)
	{
		factory := factoryFor(servicekit.MakeRegisterEngineEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.RegisterEngineEndpoint = retry
	}
	{
		factory := factoryFor(servicekit.MakeGetRegisteredEngineEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetRegisteredEngineEndpoint = retry
	}

	return endpoints, nil
}

func factoryFor(makeEndpoint func(servicekit.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := servicekit.MakeClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}