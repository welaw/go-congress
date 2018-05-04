package transport

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/welaw/go-congress/endpoints"
	"github.com/welaw/go-congress/pkg/errs"
	"github.com/welaw/go-congress/proto"
	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	sendVote grpctransport.Handler
	sendLaw  grpctransport.Handler
	status   grpctransport.Handler
}

func MakeGrpcServer(set endpoints.Endpoints, tracer stdopentracing.Tracer, logger log.Logger) proto.GoCongressServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		sendVote: grpctransport.NewServer(
			set.SendVoteEndpoint,
			DecodeGrpcSendVoteRequest,
			EncodeGrpcSendVoteResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "send_vote", logger)))...,
		),
		sendLaw: grpctransport.NewServer(
			set.SendLawEndpoint,
			DecodeGrpcSendLawRequest,
			EncodeGrpcSendLawResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "send_law", logger)))...,
		),
		status: grpctransport.NewServer(
			set.StatusEndpoint,
			DecodeGrpcStatusRequest,
			EncodeGrpcStatusResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "status", logger)))...,
		),
	}
}

func (s *grpcServer) SendVote(ctx oldcontext.Context, req *proto.SendVoteRequest) (*proto.SendVoteReply, error) {
	_, rep, err := s.sendVote.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.SendVoteReply), nil
}

func (s *grpcServer) SendLaw(ctx oldcontext.Context, req *proto.SendLawRequest) (*proto.SendLawReply, error) {
	_, rep, err := s.sendLaw.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.SendLawReply), nil
}

func (s *grpcServer) Status(ctx oldcontext.Context, req *proto.StatusRequest) (*proto.StatusReply, error) {
	_, rep, err := s.status.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.StatusReply), nil
}

func EncodeGrpcSendVoteRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.SendVoteRequest)
	return &proto.SendVoteRequest{ItemRange: req.ItemRange}, nil
}

func DecodeGrpcSendVoteRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SendVoteRequest)
	return endpoints.SendVoteRequest{ItemRange: req.ItemRange}, nil
}

func EncodeGrpcSendVoteResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.SendVoteResponse)
	return &proto.SendVoteReply{
		NewItems: resp.NewItems,
		Updated:  resp.Updated,
		Err:      errs.ErrToStr(resp.Err),
	}, nil
}

func DecodeGrpcSendVoteResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SendVoteReply)
	return endpoints.SendVoteResponse{
		NewItems: reply.NewItems,
		Updated:  reply.Updated,
		Err:      errs.StrToErr(reply.Err),
	}, nil
}

func EncodeGrpcSendLawRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.SendLawRequest)
	return &proto.SendLawRequest{ItemRange: req.ItemRange}, nil
}

func DecodeGrpcSendLawRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SendLawRequest)
	return endpoints.SendLawRequest{ItemRange: req.ItemRange}, nil
}

func EncodeGrpcSendLawResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.SendLawResponse)
	return &proto.SendLawReply{
		NewItems: resp.NewItems,
		Updated:  resp.Updated,
		Err:      errs.ErrToStr(resp.Err),
	}, nil
}

func DecodeGrpcSendLawResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.SendLawReply)
	return endpoints.SendLawResponse{
		NewItems: reply.NewItems,
		Updated:  reply.Updated,
		Err:      errs.StrToErr(reply.Err),
	}, nil
}

func EncodeGrpcStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.StatusRequest)
	return &proto.StatusRequest{ItemRange: req.ItemRange}, nil
}

func DecodeGrpcStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.StatusRequest)
	return endpoints.StatusRequest{ItemRange: req.ItemRange}, nil
}

func EncodeGrpcStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.StatusResponse)
	return &proto.StatusReply{
		NewItems: resp.NewItems,
		Existing: resp.Existing,
		Err:      errs.ErrToStr(resp.Err),
	}, nil
}

func DecodeGrpcStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.StatusReply)
	return endpoints.StatusResponse{
		NewItems: reply.NewItems,
		Existing: reply.Existing,
		Err:      errs.StrToErr(reply.Err),
	}, nil
}
