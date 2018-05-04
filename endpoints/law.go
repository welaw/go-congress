package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/welaw/go-congress/proto"
	"github.com/welaw/go-congress/services"
)

func (e Endpoints) SendLaw(ctx context.Context, scope *proto.ItemRange) (*proto.SendLawReply, error) {
	req := SendLawRequest{ItemRange: scope}
	resp, err := e.SendLawEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(SendLawResponse)
	return &proto.SendLawReply{NewItems: r.NewItems, Updated: r.Updated}, r.Err
}

func (e Endpoints) Status(ctx context.Context, scope *proto.ItemRange) (*proto.StatusReply, error) {
	req := StatusRequest{ItemRange: scope}
	resp, err := e.StatusEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	r := resp.(StatusResponse)
	return &proto.StatusReply{NewItems: r.NewItems, Existing: r.Existing}, r.Err
}

func MakeSendLawEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendLawRequest)
		resp, err := svc.SendLaw(ctx, req.ItemRange)
		if resp == nil {
			return SendLawResponse{NewItems: []string{}, Updated: []string{}, Err: err}, nil
		}
		return SendLawResponse{NewItems: resp.NewItems, Updated: resp.Updated, Err: err}, nil
	}
}

func MakeStatusEndpoint(svc services.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StatusRequest)
		res, err := svc.Status(ctx, req.ItemRange)
		return StatusResponse{
			NewItems: res.NewItems,
			Existing: res.Existing,
			Err:      err,
		}, nil
	}
}

type SendLawRequest struct {
	ItemRange *proto.ItemRange `json:"item_range"`
}

type SendLawResponse struct {
	NewItems []string `json:"new_items"`
	Updated  []string `json:"updated"`
	Err      error    `json:"err"`
}

func (r SendLawResponse) Failed() error { return r.Err }

type StatusRequest struct {
	ItemRange *proto.ItemRange `json:"item_range"`
}

type StatusResponse struct {
	NewItems []string `json:"new_items"`
	Existing []string `json:"existing"`
	Err      error    `json:"err"`
}

func (r StatusResponse) Failed() error { return r.Err }
