package services

import (
	"context"
	"fmt"

	"github.com/welaw/go-congress/congress"
	"github.com/welaw/go-congress/fdsys"
	"github.com/welaw/go-congress/pkg/errs"
	"github.com/welaw/go-congress/proto"
	welawproto "github.com/welaw/welaw/proto"
)

func (svc *service) SendVote(ctx context.Context, req *proto.ItemRange) (*proto.SendVoteReply, error) {
	var limit int32
	if req.Limit > 0 {
		limit = req.Limit
	} else {
		limit = 1
	}
	var laws []*fdsys.Item
	var err error
	if req.Ident == "" {
		laws, err = svc.db.ListLaws()
		if err != nil {
			return nil, err
		}
	} else {
		law, err := svc.db.GetItemByIdent(req.Ident)
		if err != nil {
			return nil, err
		}
		laws = []*fdsys.Item{law}
	}
	var count int32
	var resp []string
	for _, law := range laws {
		if count >= limit {
			break
		}
		c, t, version, _ := fdsys.ParseLoc(law.Loc)
		ident := fdsys.ToIdent(c, t, version)
		var chamber string
		if t == congress.HouseKey {
			chamber = congress.HouseIdent
		} else {
			chamber = congress.SenateIdent
		}

		err = svc.db.CheckRollCall(ident)
		switch {
		case err == errs.ErrNotFound:
		case err != nil:
			return nil, err
		default:
			svc.logger.Log("method", "send_vote", "roll call found", "skipping")
			continue
		}
		err = svc.sendRollCall(c, chamber, "1", version, ident)
		if err != nil {
			return nil, err
		}
		err = svc.db.CreateRollCall(ident)
		if err != nil {
			return nil, err
		}
		count++
		resp = append(resp, fmt.Sprintf("%v", count))
	}

	return &proto.SendVoteReply{NewItems: resp}, nil
}

func (svc *service) sendRollCall(c, chamber, session, billNumber, ident string) error {
	svc.logger.Log("method", "send_roll_call", "congress", c, "chamber", chamber, "session", session, "billNumber", billNumber)

	rollCall, result, err := congress.GetRollCalls(c, chamber, session, billNumber)
	if err != nil {
		return err
	}
	if len(rollCall) == 0 {
		svc.logger.Log("warning", "no votes found", "congress", c, "chamber", chamber, "session", session, "bill_number", billNumber)
		return nil
	}

	var user *proto.User
	users := make(map[string]*proto.User)
	ctx := context.Background()
	var votes []*welawproto.Vote
	var newUsers []string
	for id := range rollCall {
		// lisRegex := regexp.MustCompile(`^[a-zA-Z]{1}\d{3}$`)
		// bioguideRegex := regexp.MustCompile(`^[a-zA-Z]{1}\d{6}$`)
		var err error
		switch {
		case lisRegex.MatchString(id):
			user, err = svc.db.GetUserByLisId(id)
		case bioguideRegex.MatchString(id):
			user, err = svc.db.GetUserByBioguideId(id)
		default:
			return fmt.Errorf("unknown id: %v", id)
		}
		if err == errs.ErrNotFound {
			svc.logger.Log("user not found in local db, creating", id)
			newUsers = append(newUsers, id)
		} else if err != nil {
			return err
		} else {
			//color.Yellow("adding existing user: %v", user)
			users[id] = user
		}
	}

	output := make(chan *proto.User, 1000)
	input := make(chan string, 1000)
	errc := make(chan error, 1000)

	for w := 0; w < 10; w++ {
		go svc.getUserWorker(ctx, input, output, errc)
	}
	for _, id := range newUsers {
		input <- id
	}
	for i := 0; i < len(newUsers); i++ {
		//color.Yellow("waiting for new user worker")
		select {
		case err := <-errc:
			//errs = append(errs, err)
			//break
			panic(err)
		case user := <-output:
			if user == nil {
				svc.logger.Log("user not found", newUsers[i])
				continue
			}
			_, err := svc.db.GetUserByBioguideId(user.BioguideId)
			if err == errs.ErrNotFound {
			} else if err != nil {
				return err
			} else {
				//color.Yellow("new user found in database")
				continue
			}
			users[user.FoundId] = user
		}
	}

	// for _, id := range newUsers {
	// 	u, err := svc.getUser(id)
	// 	if err == errs.ErrNotFound {
	// 		// TODO
	// 		// remove votes by user
	// 		delete(rollCall, id)
	// 		svc.logger.Log("user not found, skipping", id)
	// 		continue
	// 	} else if err != nil {
	// 		return err
	// 	}
	// 	color.Yellow("appending new user: %v", u)
	// 	users[id] = u
	// }

	for id, v := range rollCall {
		user, ok := users[id]
		if !ok {
			return fmt.Errorf("user not found in map")
		}
		if user.Username == "" {
			panic(fmt.Errorf("username not found: %v, %v", id, v))
		}
		vote := &welawproto.Vote{
			Vote:     v,
			Username: user.Username,
			Upstream: upstreamIdent,
			Ident:    ident,
			Branch:   "master",
		}
		_, err = svc.db.GetVote(vote)
		switch {
		case err == errs.ErrNotFound:
		case err != nil:
			return err
		default:
			svc.logger.Log("vote exists in db, continuing...", fmt.Sprintf("%+v", vote))
			continue
		}

		//color.Yellow("adding new vote: %v", vote)
		votes = append(votes, vote)
	}

	var wusers []*welawproto.User
	for _, u := range users {
		wusers = append(wusers, temp(u))
	}

	//color.Yellow("sending create users: %v", len(wusers))
	_, err = svc.client.CreateUsers(ctx, wusers)
	for _, u := range users {
		err = svc.db.CreateUser(u)
		if err != nil {
			return err
		}
	}

	//color.Yellow("sending create votes: %v", len(votes))
	opts := &welawproto.CreateVotesOptions{
		VoteResult: result,
	}
	_, err = svc.client.CreateVotes(ctx, votes, opts)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) getUserWorker(ctx context.Context, input <-chan string, output chan<- *proto.User, errc chan<- error) {
	for id := range input {
		//color.Yellow("ensure user worker, input")
		u, err := svc.getUser(id)
		if err == errs.ErrNotFound {
			output <- nil
			continue
		} else if err != nil {
			panic(fmt.Errorf("error ensuring user group --- %v", err))
			//errc <- err
			//return
		}
		u.FoundId = id
		output <- u
	}
}
