// Package model contains the types for schema 'public'.
package model

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// TxTransfer represents a row from 'public.tx_transfer'.
type TxTransfer struct {
	ID      int64          `json:"id"`       // id
	ChainID sql.NullString `json:"chain_id"` // chain_id
	Height  sql.NullInt64  `json:"height"`   // height
	Hash    sql.NullString `json:"hash"`     // hash
	Address sql.NullString `json:"address"`  // address
	Coin    sql.NullString `json:"coin"`     // coin
	Amount  sql.NullString `json:"amount"`   // amount
	Type    sql.NullInt64  `json:"type"`     // type
	Time    pq.NullTime    `json:"time"`     // time

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the TxTransfer exists in the database.
func (tt *TxTransfer) Exists() bool {
	return tt._exists
}

// Deleted provides information if the TxTransfer has been deleted from the database.
func (tt *TxTransfer) Deleted() bool {
	return tt._deleted
}

// Insert inserts the TxTransfer to the database.
func (tt *TxTransfer) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if tt._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.tx_transfer (` +
		`chain_id, height, hash, address, coin, amount, type, time` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8` +
		`) RETURNING id`

	// run query
	XOLog(sqlstr, tt.ChainID, tt.Height, tt.Hash, tt.Address, tt.Coin, tt.Amount, tt.Type, tt.Time)
	err = db.QueryRow(sqlstr, tt.ChainID, tt.Height, tt.Hash, tt.Address, tt.Coin, tt.Amount, tt.Type, tt.Time).Scan(&tt.ID)
	if err != nil {
		return err
	}

	// set existence
	tt._exists = true

	return nil
}

// Update updates the TxTransfer in the database.
func (tt *TxTransfer) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !tt._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if tt._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.tx_transfer SET (` +
		`chain_id, height, hash, address, coin, amount, type, time` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7, $8` +
		`) WHERE id = $9`

	// run query
	XOLog(sqlstr, tt.ChainID, tt.Height, tt.Hash, tt.Address, tt.Coin, tt.Amount, tt.Type, tt.Time, tt.ID)
	_, err = db.Exec(sqlstr, tt.ChainID, tt.Height, tt.Hash, tt.Address, tt.Coin, tt.Amount, tt.Type, tt.Time, tt.ID)
	return err
}

// Save saves the TxTransfer to the database.
func (tt *TxTransfer) Save(db XODB) error {
	if tt.Exists() {
		return tt.Update(db)
	}

	return tt.Insert(db)
}

// Upsert performs an upsert for TxTransfer.
//
// NOTE: PostgreSQL 9.5+ only
func (tt *TxTransfer) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if tt._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.tx_transfer (` +
		`id, chain_id, height, hash, address, coin, amount, type, time` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, chain_id, height, hash, address, coin, amount, type, time` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.chain_id, EXCLUDED.height, EXCLUDED.hash, EXCLUDED.address, EXCLUDED.coin, EXCLUDED.amount, EXCLUDED.type, EXCLUDED.time` +
		`)`

	// run query
	XOLog(sqlstr, tt.ID, tt.ChainID, tt.Height, tt.Hash, tt.Address, tt.Coin, tt.Amount, tt.Type, tt.Time)
	_, err = db.Exec(sqlstr, tt.ID, tt.ChainID, tt.Height, tt.Hash, tt.Address, tt.Coin, tt.Amount, tt.Type, tt.Time)
	if err != nil {
		return err
	}

	// set existence
	tt._exists = true

	return nil
}

// Delete deletes the TxTransfer from the database.
func (tt *TxTransfer) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !tt._exists {
		return nil
	}

	// if deleted, bail
	if tt._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.tx_transfer WHERE id = $1`

	// run query
	XOLog(sqlstr, tt.ID)
	_, err = db.Exec(sqlstr, tt.ID)
	if err != nil {
		return err
	}

	// set deleted
	tt._deleted = true

	return nil
}

// TxTransfersQuery returns offset-limit rows from 'public.tx_transfer' filte by filter,
// ordered by "id" in descending order.
func TxTransferFilter(db XODB, filter, sort string, offset, limit int64) ([]*TxTransfer, error) {
	sqlstr := `SELECT ` +
		`id, chain_id, height, hash, address, coin, amount, type, time` +
		` FROM public.tx_transfer `

	if filter != "" {
		sqlstr = sqlstr + " WHERE " + filter
	}

	if sort != "" {
		sqlstr = sqlstr + " " + sort
	}

	if limit > 0 {
		sqlstr = sqlstr + fmt.Sprintf(" offset %d limit %d", offset, limit)
	}

	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*TxTransfer
	for q.Next() {
		tt := TxTransfer{
			_exists: true,
		}

		// scan
		err = q.Scan(&tt.ID, &tt.ChainID, &tt.Height, &tt.Hash, &tt.Address, &tt.Coin, &tt.Amount, &tt.Type, &tt.Time)
		if err != nil {
			return nil, err
		}

		res = append(res, &tt)
	}

	return res, nil
}

// TxTransfersByAddress retrieves a row from 'public.tx_transfer' as a TxTransfer.
//
// Generated from index 'tx_transfer_address_idx'.
func TxTransfersByAddress(db XODB, address sql.NullString) ([]*TxTransfer, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, chain_id, height, hash, address, coin, amount, type, time ` +
		`FROM public.tx_transfer ` +
		`WHERE address = $1`

	// run query
	XOLog(sqlstr, address)
	q, err := db.Query(sqlstr, address)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*TxTransfer{}
	for q.Next() {
		tt := TxTransfer{
			_exists: true,
		}

		// scan
		err = q.Scan(&tt.ID, &tt.ChainID, &tt.Height, &tt.Hash, &tt.Address, &tt.Coin, &tt.Amount, &tt.Type, &tt.Time)
		if err != nil {
			return nil, err
		}

		res = append(res, &tt)
	}

	return res, nil
}

// TxTransferByID retrieves a row from 'public.tx_transfer' as a TxTransfer.
//
// Generated from index 'tx_transfer_pkey'.
func TxTransferByID(db XODB, id int64) (*TxTransfer, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, chain_id, height, hash, address, coin, amount, type, time ` +
		`FROM public.tx_transfer ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	tt := TxTransfer{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&tt.ID, &tt.ChainID, &tt.Height, &tt.Hash, &tt.Address, &tt.Coin, &tt.Amount, &tt.Type, &tt.Time)
	if err != nil {
		return nil, err
	}

	return &tt, nil
}
