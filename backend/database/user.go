package database

import (
	"database/sql"

	"github.com/welaw/go-congress/pkg/errs"
	"github.com/welaw/go-congress/proto"
)

func (db *_database) CreateUser(user *proto.User) error {
	q := `
	INSERT INTO users (
		username,
		bioguide_id,
		lis_id
	)
	VALUES ($1, $2, $3)
	`
	_, err := db.conn.Exec(q, user.Username, user.BioguideId, user.LisId)
	return err
}

func (db *_database) GetUserByBioguideId(bioguideId string) (*proto.User, error) {
	//color.Blue("getting by bioguide db: %v", bioguideId)
	q := `
	SELECT username,
		lis_id
	FROM users
	WHERE bioguide_id = $1
		AND	deleted_at IS NULL
	`
	var u proto.User
	err := db.conn.QueryRow(q, bioguideId).Scan(&u.Username, &u.LisId)
	if err == sql.ErrNoRows {
		return nil, errs.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	u.BioguideId = bioguideId
	//color.Blue("user: %v", u)
	return &u, err
}

func (db *_database) GetUserByLisId(lisId string) (*proto.User, error) {
	//color.Blue("getting by lis_id: %v", lisId)
	q := `
	SELECT username,
		bioguide_id
	FROM users
	WHERE lis_id = $1
		AND	deleted_at IS NULL
	`
	var u proto.User
	err := db.conn.QueryRow(q, lisId).Scan(&u.Username, &u.BioguideId)
	if err == sql.ErrNoRows {
		return nil, errs.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	u.LisId = lisId
	//color.Blue("user: %v", u)
	return &u, nil
}

func (db *_database) UpdateUser(user *proto.User) error {
	q := `
	UPDATE users
	SET username = (
		SELECT CASE WHEN $1 = '' THEN username
			ELSE $1
		END
		FROM users WHERE 
	), bioguide_id = (
	), lis_id = (
		SELECT
	)
	VALUES ($1, $2, $3)
	`
	_, err := db.conn.Exec(q, user.Username, user.BioguideId, user.LisId)
	return err
}
