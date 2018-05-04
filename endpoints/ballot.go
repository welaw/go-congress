package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/welaw/go-congress/proto"
	"github.com/welaw/go-congress/services"
)

func (e Endpoints) SendVote(ctx context.Context, scope *proto.ItemRange) (*proto.SendVoteReply, error) {
	req := SendVoteRequest{ItemRange: scope}
	resp, err := e.SendVoteEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(SendVoteResponse)
	return &proto.SendVoteReply{NewItems: r.NewItems, Updated: r.Updated}, r.Err
}

func MakeSendVoteEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendVoteRequest)
		resp, err := svc.SendVote(ctx, req.ItemRange)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return SendVoteResponse{NewItems: []string{}, Updated: []string{}}, nil
		}
		return SendVoteResponse{NewItems: resp.NewItems, Updated: resp.Updated, Err: err}, nil
	}
}

type SendVoteRequest struct {
	ItemRange *proto.ItemRange `json:"item_range"`
}

type SendVoteResponse struct {
	NewItems []string `json:"new_items"`
	Updated  []string `json:"updated"`
	Err      error    `json:"err"`
}

func (r SendVoteResponse) Failed() error { return r.Err }
