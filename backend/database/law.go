package database

import (
	"database/sql"

	"github.com/welaw/go-congress/fdsys"
	"github.com/welaw/go-congress/pkg/errs"
)

func (db *_database) CreateItem(item *fdsys.Item) error {
	q := `
	INSERT INTO items (
		ident,
		bioguide_id,
		loc,
		last_mod
	)
	VALUES ($1, $2, $3, $4)
	RETURNING items.loc
	`
	var dummy string
	err := db.conn.QueryRow(
		q,
		item.Ident,
		item.BioguideID,
		item.Loc,
		item.LastMod,
	).Scan(&dummy)
	if err != nil {
		return err
	}
	return nil
}

func (db *_database) GetItemByIdent(ident string) (*fdsys.Item, error) {
	//color.Blue("getting by ident db: %v", ident)
	q := `
	SELECT loc,
		last_mod
	FROM items
	WHERE items.ident = $1
	AND	deleted_at IS NULL
	`
	var item fdsys.Item
	err := db.conn.QueryRow(q, ident).Scan(
		&item.Loc,
		&item.LastMod,
	)
	if err == sql.ErrNoRows {
		return nil, errs.ErrNotFound
	}
	return &item, err
}

func (db *_database) GetItemByLoc(loc string) (*fdsys.Item, error) {
	q := `
	SELECT ident,
		last_mod
	FROM items
	WHERE items.loc = $1
	AND	deleted_at IS NULL
	`
	var item fdsys.Item
	err := db.conn.QueryRow(q, loc).Scan(
		&item.Ident,
		&item.LastMod,
	)
	if err == sql.ErrNoRows {
		return nil, errs.ErrNotFound
	}
	return &item, err
}

func (db *_database) GetPreviousVersion(ident string) (*fdsys.Item, error) {
	q := `
	SELECT ident,
		loc,
		last_mod
	FROM items
	WHERE position($1 IN items.loc) <> 0
	AND	deleted_at IS NULL
	ORDER BY last_mod DESC
	RETURNING 1
	`
	var item fdsys.Item
	err := db.conn.QueryRow(q, ident).Scan(
		&item.Ident,
		&item.Loc,
		&item.LastMod,
	)
	return &item, err
}

func (db *_database) GetLawItems(ident string) ([]*fdsys.Item, error) {
	q := `
	SELECT ident,
		loc,
		last_mod
	FROM items
	WHERE position($1 IN items.loc) <> 0
	AND	deleted_at IS NULL
	ORDER BY last_mod
	`
	var items []*fdsys.Item
	rows, err := db.conn.Query(q, ident)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i fdsys.Item
		if err = rows.Scan(
			&i.Ident,
			&i.Loc,
			&i.LastMod,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	err = rows.Err()
	return items, err
}

//func (db *_database) ListItems(max int32) ([]*fdsys.Item, error) {
func (db *_database) ListItems() ([]*fdsys.Item, error) {
	q := `
	SELECT ident,
		loc,
		last_mod
	FROM items
	WHERE deleted_at IS NULL
	`
	var items []*fdsys.Item
	rows, err := db.conn.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i fdsys.Item
		if err = rows.Scan(
			&i.Ident,
			&i.Loc,
			&i.LastMod,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	err = rows.Err()
	return items, err
}

func (db *_database) ListLaws() ([]*fdsys.Item, error) {
	q := `
	SELECT DISTINCT ON (ident)
		ident,
		loc,
		last_mod
	FROM items
	WHERE deleted_at IS NULL
	`
	var items []*fdsys.Item
	rows, err := db.conn.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i fdsys.Item
		if err = rows.Scan(
			&i.Ident,
			&i.Loc,
			&i.LastMod,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	err = rows.Err()
	return items, err
}

func (db *_database) UpdateItem(newItem *fdsys.Item) error {
	// get current
	item, err := db.GetItemByLoc(newItem.Loc)
	if err != nil {
		return err
	}

	copyItem(item, newItem)

	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	// delete old
	delete := `
	UPDATE items
	SET deleted_at = $2
	WHERE loc = $1
	AND deleted_at IS NULL
	RETURNING 1
	`
	stmt, err := tx.Prepare(delete)
	//rows, err := tx.Query(
	res, err := stmt.Exec(
		item.Loc,
		item.LastMod,
	)
	rows, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows == 0 {
		tx.Rollback()
		return err
	}

	// create new
	create := `
	INSERT INTO items (
		ident,
		loc,
		last_mod
	)
	VALUES ($1, $2)
	`
	_, err = tx.Exec(
		create,
		item.Ident,
		item.Loc,
		item.LastMod,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func copyItem(from, to *fdsys.Item) {
	if from.Ident != "" {
		to.Ident = from.Ident
	}
	if from.Loc != "" {
		to.Loc = from.Loc
	}
	if from.LastMod != "" {
		to.LastMod = from.LastMod
	}
}
