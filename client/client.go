package client

import (
	"context"

	"github.com/welaw/go-congress/proto"
)

type GoCongressClient interface {
	// ballot
	SendVote(context.Context, *proto.ItemRange) (*proto.SendVoteReply, error)
	// law
	SendLaw(context.Context, *proto.ItemRange) (*proto.SendLawReply, error)
	Status(context.Context, *proto.ItemRange) (*proto.StatusReply, error)
}
