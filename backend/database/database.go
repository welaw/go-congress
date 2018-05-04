package database

import (
	"database/sql"

	"github.com/welaw/go-congress/fdsys"
	"github.com/welaw/go-congress/proto"
	welawproto "github.com/welaw/welaw/proto"
)

type Database interface {
	// ballot
	CreateVote(*welawproto.Vote) error
	GetVote(*welawproto.Vote) (*welawproto.Vote, error)
	CreateRollCall(string) error
	CheckRollCall(string) error
	// law
	CreateItem(*fdsys.Item) error
	GetItemByIdent(string) (*fdsys.Item, error)
	GetItemByLoc(string) (*fdsys.Item, error)
	GetLawItems(string) ([]*fdsys.Item, error)
	ListItems() ([]*fdsys.Item, error)
	ListLaws() ([]*fdsys.Item, error)
	UpdateItem(*fdsys.Item) error
	// user
	CreateUser(*proto.User) error
	GetUserByBioguideId(string) (*proto.User, error)
	GetUserByLisId(string) (*proto.User, error)

	Close() error
}

type _database struct {
	conn *sql.DB
}

func (db *_database) Close() error {
	return db.conn.Close()
}

func openDatabase(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDatabase(url string) (Database, error) {
	db, err := openDatabase(url)
	if err != nil {
		return nil, err
	}
	return &_database{
		conn: db,
	}, nil
}
