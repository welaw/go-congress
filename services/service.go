package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/welaw/go-congress/backend/database"
	"github.com/welaw/go-congress/backend/filesystem"
	"github.com/welaw/go-congress/pkg/errs"
	"github.com/welaw/go-congress/pkg/utils"
	"github.com/welaw/go-congress/proto"
	welawproto "github.com/welaw/welaw/proto"
	"github.com/welaw/welaw/services"
)

const (
	upstreamIdent      = "congress"
	upstreamProvider   = "welaw"
	upstreamPictureURL = "/assets/congress.png"
)

type Service interface {
	// ballot
	SendVote(context.Context, *proto.ItemRange) (*proto.SendVoteReply, error)
	// law
	SendLaw(context.Context, *proto.ItemRange) (*proto.SendLawReply, error)
	Status(context.Context, *proto.ItemRange) (*proto.StatusReply, error)
}

type service struct {
	db     database.Database
	fs     filesystem.Filesystem
	logger log.Logger
	// TODO
	upstream welawproto.Upstream
	client   services.Service
}

func NewService(
	db database.Database,
	fs filesystem.Filesystem,
	logger log.Logger,
	upstream welawproto.Upstream,
	client services.Service,
) Service {
	return &service{
		db:       db,
		fs:       fs,
		logger:   logger,
		upstream: upstream,
		client:   client,
	}
}

func (svc *service) getFile(url string) ([]byte, error) {
	s := strings.Split(url, "/")
	var filename string
	switch {
	case len(s) > 1:
		filename = strings.Join(s[len(s)-2:], "-")
	case len(s) == 1:
		filename = s[0]
	default:
		return nil, fmt.Errorf("bad url: %v", url)
	}
	f, err := svc.fs.OpenFile(filename)
	switch {
	case err == errs.ErrNotFound:
		svc.logger.Log("method", "get_file", "file not found, fetching remote", url)
		f, err = utils.DownloadData(url)
		if err != nil {
			return nil, err
		}
		err = svc.fs.WriteFile(filename, f)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	default:
		svc.logger.Log("method", "get_file", "using local file", filename)
	}
	return f, nil
}
