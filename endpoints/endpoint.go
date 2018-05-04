package endpoints

import (
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	// ballot
	SendVoteEndpoint endpoint.Endpoint
	// law
	SendLawEndpoint endpoint.Endpoint
	StatusEndpoint  endpoint.Endpoint
}

type Failer interface {
	Failed() error
}
