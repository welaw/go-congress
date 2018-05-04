package database

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/welaw/go-congress/pkg/errs"
	welawproto "github.com/welaw/welaw/proto"
)

func (db *_database) CreateVote(vote *welawproto.Vote) (err error) {
	q := `
	INSERT INTO votes (
		username,
		vote,
		upstream,
		ident,
		branch
	)
	VALUES ($1, $2, $3, $4, $5)
	`
	err = db.conn.QueryRow(
		q,
		vote.Username,
		vote.Vote,
		vote.Upstream,
		vote.Ident,
		vote.Branch,
	).Scan()
	if err != nil {
		return err
	}
	return nil
}

func (db *_database) GetVote(vote *welawproto.Vote) (v *welawproto.Vote, err error) {
	q := `
	SELECT uid
	FROM votes
	WHERE votes.username = $1
		AND votes.vote = $2
		AND votes.upstream = $3
		AND votes.ident = $4
		AND votes.branch = $5
		AND	deleted_at IS NULL
	`
	var uid uuid.UUID
	err = db.conn.QueryRow(
		q,
		vote.Username,
		vote.Vote,
		vote.Upstream,
		vote.Ident,
		vote.Branch,
	).Scan(&uid)
	switch {
	case err == sql.ErrNoRows:
		return nil, errs.ErrNotFound
	case err != nil:
		return nil, err
	}
	vote.Uid = uid.String()
	return vote, nil
}

func (db *_database) CreateRollCall(ident string) (err error) {
	q := `
	INSERT INTO rollcall (
		ident
	)
	VALUES ($1)
	`
	_, err = db.conn.Exec(q, ident)
	if err != nil {
		return err
	}
	return nil
}

func (db *_database) CheckRollCall(ident string) (err error) {
	q := `
	SELECT 1
	FROM rollcall
	WHERE ident = $1
	AND	deleted_at IS NULL
	`
	var d int
	err = db.conn.QueryRow(
		q,
		ident,
	).Scan(&d)
	switch {
	case err == sql.ErrNoRows:
		return errs.ErrNotFound
	case err != nil:
		return err
	}
	return nil
}
