// Package model contains the types for schema 'public'.
package model

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

// Tx represents a row from 'public.txs'.
type Tx struct {
	ID          int64          `json:"id"`           // id
	ChainID     sql.NullString `json:"chain_id"`     // chain_id
	Height      sql.NullInt64  `json:"height"`       // height
	TxType      sql.NullString `json:"tx_type"`      // tx_type
	Index       sql.NullInt64  `json:"index"`        // index
	Maxgas      sql.NullInt64  `json:"maxgas"`       // maxgas
	QcpFrom     sql.NullString `json:"qcp_from"`     // qcp_from
	QcpTo       sql.NullString `json:"qcp_to"`       // qcp_to
	QcpSequence sql.NullInt64  `json:"qcp_sequence"` // qcp_sequence
	QcpTxindex  sql.NullInt64  `json:"qcp_txindex"`  // qcp_txindex
	QcpIsresult sql.NullBool   `json:"qcp_isresult"` // qcp_isresult
	OriginTx    sql.NullString `json:"origin_tx"`    // origin_tx
	JSONTx      sql.NullString `json:"json_tx"`      // json_tx
	Time        pq.NullTime    `json:"time"`         // time
	CreatedAt   pq.NullTime    `json:"created_at"`   // created_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Tx exists in the database.
func (t *Tx) Exists() bool {
	return t._exists
}

// Deleted provides information if the Tx has been deleted from the database.
func (t *Tx) Deleted() bool {
	return t._deleted
}

// Insert inserts the Tx to the database.
func (t *Tx) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if t._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.txs (` +
		`chain_id, height, tx_type, index, maxgas, qcp_from, qcp_to, qcp_sequence, qcp_txindex, qcp_isresult, origin_tx, json_tx, time, created_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14` +
		`) RETURNING id`

	// run query
	XOLog(sqlstr, t.ChainID, t.Height, t.TxType, t.Index, t.Maxgas, t.QcpFrom, t.QcpTo, t.QcpSequence, t.QcpTxindex, t.QcpIsresult, t.OriginTx, t.JSONTx, t.Time, t.CreatedAt)
	err = db.QueryRow(sqlstr, t.ChainID, t.Height, t.TxType, t.Index, t.Maxgas, t.QcpFrom, t.QcpTo, t.QcpSequence, t.QcpTxindex, t.QcpIsresult, t.OriginTx, t.JSONTx, t.Time, t.CreatedAt).Scan(&t.ID)
	if err != nil {
		return err
	}

	// set existence
	t._exists = true

	return nil
}

// Update updates the Tx in the database.
func (t *Tx) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !t._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if t._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.txs SET (` +
		`chain_id, height, tx_type, index, maxgas, qcp_from, qcp_to, qcp_sequence, qcp_txindex, qcp_isresult, origin_tx, json_tx, time, created_at` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14` +
		`) WHERE id = $15`

	// run query
	XOLog(sqlstr, t.ChainID, t.Height, t.TxType, t.Index, t.Maxgas, t.QcpFrom, t.QcpTo, t.QcpSequence, t.QcpTxindex, t.QcpIsresult, t.OriginTx, t.JSONTx, t.Time, t.CreatedAt, t.ID)
	_, err = db.Exec(sqlstr, t.ChainID, t.Height, t.TxType, t.Index, t.Maxgas, t.QcpFrom, t.QcpTo, t.QcpSequence, t.QcpTxindex, t.QcpIsresult, t.OriginTx, t.JSONTx, t.Time, t.CreatedAt, t.ID)
	return err
}

// Save saves the Tx to the database.
func (t *Tx) Save(db XODB) error {
	if t.Exists() {
		return t.Update(db)
	}

	return t.Insert(db)
}

// Upsert performs an upsert for Tx.
//
// NOTE: PostgreSQL 9.5+ only
func (t *Tx) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if t._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.txs (` +
		`id, chain_id, height, tx_type, index, maxgas, qcp_from, qcp_to, qcp_sequence, qcp_txindex, qcp_isresult, origin_tx, json_tx, time, created_at` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, chain_id, height, tx_type, index, maxgas, qcp_from, qcp_to, qcp_sequence, qcp_txindex, qcp_isresult, origin_tx, json_tx, time, created_at` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.chain_id, EXCLUDED.height, EXCLUDED.tx_type, EXCLUDED.index, EXCLUDED.maxgas, EXCLUDED.qcp_from, EXCLUDED.qcp_to, EXCLUDED.qcp_sequence, EXCLUDED.qcp_txindex, EXCLUDED.qcp_isresult, EXCLUDED.origin_tx, EXCLUDED.json_tx, EXCLUDED.time, EXCLUDED.created_at` +
		`)`

	// run query
	XOLog(sqlstr, t.ID, t.ChainID, t.Height, t.TxType, t.Index, t.Maxgas, t.QcpFrom, t.QcpTo, t.QcpSequence, t.QcpTxindex, t.QcpIsresult, t.OriginTx, t.JSONTx, t.Time, t.CreatedAt)
	_, err = db.Exec(sqlstr, t.ID, t.ChainID, t.Height, t.TxType, t.Index, t.Maxgas, t.QcpFrom, t.QcpTo, t.QcpSequence, t.QcpTxindex, t.QcpIsresult, t.OriginTx, t.JSONTx, t.Time, t.CreatedAt)
	if err != nil {
		return err
	}

	// set existence
	t._exists = true

	return nil
}

// Delete deletes the Tx from the database.
func (t *Tx) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !t._exists {
		return nil
	}

	// if deleted, bail
	if t._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.txs WHERE id = $1`

	// run query
	XOLog(sqlstr, t.ID)
	_, err = db.Exec(sqlstr, t.ID)
	if err != nil {
		return err
	}

	// set deleted
	t._deleted = true

	return nil
}

// TxsQuery returns offset-limit rows from 'public.txs' filte by filter,
// ordered by "id" in descending order.
func TxFilter(db XODB, filter string, offset, limit int64) ([]*Tx, error) {
	sqlstr := `SELECT ` +
		`id, chain_id, height, tx_type, index, maxgas, qcp_from, qcp_to, qcp_sequence, qcp_txindex, qcp_isresult, origin_tx, json_tx, time, created_at` +
		` FROM public.txs `

	if filter != "" {
		sqlstr = sqlstr + " WHERE " + filter
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
	var res []*Tx
	for q.Next() {
		t := Tx{}

		// scan
		err = q.Scan(&t.ID, &t.ChainID, &t.Height, &t.TxType, &t.Index, &t.Maxgas, &t.QcpFrom, &t.QcpTo, &t.QcpSequence, &t.QcpTxindex, &t.QcpIsresult, &t.OriginTx, &t.JSONTx, &t.Time, &t.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}

// TxesByChainIDHeight retrieves a row from 'public.txs' as a Tx.
//
// Generated from index 'txs_chain_id_height_idx'.
func TxesByChainIDHeight(db XODB, chainID sql.NullString, height sql.NullInt64) ([]*Tx, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, chain_id, height, tx_type, index, maxgas, qcp_from, qcp_to, qcp_sequence, qcp_txindex, qcp_isresult, origin_tx, json_tx, time, created_at ` +
		`FROM public.txs ` +
		`WHERE chain_id = $1 AND height = $2`

	// run query
	XOLog(sqlstr, chainID, height)
	q, err := db.Query(sqlstr, chainID, height)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Tx{}
	for q.Next() {
		t := Tx{
			_exists: true,
		}

		// scan
		err = q.Scan(&t.ID, &t.ChainID, &t.Height, &t.TxType, &t.Index, &t.Maxgas, &t.QcpFrom, &t.QcpTo, &t.QcpSequence, &t.QcpTxindex, &t.QcpIsresult, &t.OriginTx, &t.JSONTx, &t.Time, &t.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}

// TxByChainIDHeightIndex retrieves a row from 'public.txs' as a Tx.
//
// Generated from index 'txs_chain_id_height_index_idx'.
func TxByChainIDHeightIndex(db XODB, chainID sql.NullString, height sql.NullInt64, index sql.NullInt64) (*Tx, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, chain_id, height, tx_type, index, maxgas, qcp_from, qcp_to, qcp_sequence, qcp_txindex, qcp_isresult, origin_tx, json_tx, time, created_at ` +
		`FROM public.txs ` +
		`WHERE chain_id = $1 AND height = $2 AND index = $3`

	// run query
	XOLog(sqlstr, chainID, height, index)
	t := Tx{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, chainID, height, index).Scan(&t.ID, &t.ChainID, &t.Height, &t.TxType, &t.Index, &t.Maxgas, &t.QcpFrom, &t.QcpTo, &t.QcpSequence, &t.QcpTxindex, &t.QcpIsresult, &t.OriginTx, &t.JSONTx, &t.Time, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// TxesByChainID retrieves a row from 'public.txs' as a Tx.
//
// Generated from index 'txs_chain_id_idx'.
func TxesByChainID(db XODB, chainID sql.NullString) ([]*Tx, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, chain_id, height, tx_type, index, maxgas, qcp_from, qcp_to, qcp_sequence, qcp_txindex, qcp_isresult, origin_tx, json_tx, time, created_at ` +
		`FROM public.txs ` +
		`WHERE chain_id = $1`

	// run query
	XOLog(sqlstr, chainID)
	q, err := db.Query(sqlstr, chainID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Tx{}
	for q.Next() {
		t := Tx{
			_exists: true,
		}

		// scan
		err = q.Scan(&t.ID, &t.ChainID, &t.Height, &t.TxType, &t.Index, &t.Maxgas, &t.QcpFrom, &t.QcpTo, &t.QcpSequence, &t.QcpTxindex, &t.QcpIsresult, &t.OriginTx, &t.JSONTx, &t.Time, &t.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}

// TxByID retrieves a row from 'public.txs' as a Tx.
//
// Generated from index 'txs_pkey'.
func TxByID(db XODB, id int64) (*Tx, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, chain_id, height, tx_type, index, maxgas, qcp_from, qcp_to, qcp_sequence, qcp_txindex, qcp_isresult, origin_tx, json_tx, time, created_at ` +
		`FROM public.txs ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	t := Tx{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&t.ID, &t.ChainID, &t.Height, &t.TxType, &t.Index, &t.Maxgas, &t.QcpFrom, &t.QcpTo, &t.QcpSequence, &t.QcpTxindex, &t.QcpIsresult, &t.OriginTx, &t.JSONTx, &t.Time, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
