package client

import (
	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/welaw/go-congress/endpoints"
	"github.com/welaw/go-congress/proto"
	"github.com/welaw/go-congress/transport"
)

func MakeGoCongressClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) GoCongressClient {

	var sendVoteEndpoint endpoint.Endpoint
	{
		sendVoteEndpoint = grpctransport.NewClient(
			conn,
			"grpc.gocongress.v1.GoCongressService",
			"SendVote",
			transport.EncodeGrpcSendVoteRequest,
			transport.DecodeGrpcSendVoteResponse,
			proto.SendVoteReply{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		sendVoteEndpoint = opentracing.TraceClient(tracer, "send_vote")(sendVoteEndpoint)
	}

	var sendLawEndpoint endpoint.Endpoint
	{
		sendLawEndpoint = grpctransport.NewClient(
			conn,
			"grpc.gocongress.v1.GoCongressService",
			"SendLaw",
			transport.EncodeGrpcSendLawRequest,
			transport.DecodeGrpcSendLawResponse,
			proto.SendLawReply{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		sendLawEndpoint = opentracing.TraceClient(tracer, "send_law")(sendLawEndpoint)
	}

	var statusEndpoint endpoint.Endpoint
	{
		statusEndpoint = grpctransport.NewClient(
			conn,
			"grpc.gocongress.v1.GoCongressService",
			"Status",
			transport.EncodeGrpcStatusRequest,
			transport.DecodeGrpcStatusResponse,
			proto.StatusReply{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		statusEndpoint = opentracing.TraceClient(tracer, "status")(statusEndpoint)
	}

	return endpoints.Endpoints{
		SendVoteEndpoint: sendVoteEndpoint,
		SendLawEndpoint:  sendLawEndpoint,
		StatusEndpoint:   statusEndpoint,
	}

}
