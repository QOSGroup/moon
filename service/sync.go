// Copyright 2018 The QOS Authors

package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/QOSGroup/qmoon/db"
	"github.com/QOSGroup/qmoon/db/model"
	"github.com/QOSGroup/qmoon/lib"
	"github.com/QOSGroup/qmoon/service/block"
	"github.com/QOSGroup/qmoon/service/tx"
	"github.com/QOSGroup/qmoon/service/validator"
	"github.com/QOSGroup/qmoon/types"
	"github.com/QOSGroup/qmoon/utils"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// CreateBlock 创建一个块
func CreateBlock(cli *lib.TmClient, height *int64) error {
	b, err := cli.RetrieveBlock(height)
	if err != nil {
		return err
	}

	err = block.Save(b)
	if err != nil {
		return err
	}

	err = tx.Save(cli, b)
	// TODO delete block

	err = validator.SaveBlockValidator(b.Precommits)
	// TODO delete block and tx

	return nil
}

func CreateValidator(chainId string, vals types.Validators) error {
	for _, val := range vals {
		validator.CreateValidator(chainId, val)
	}

	valMap := make(map[string]types.Validator)
	for _, v := range vals {
		valMap[v.Address] = v
	}

	allVals, err := validator.ListValidatorByChain(chainId)
	if err == nil {
		for _, v := range allVals {
			if v.Status == types.Active {
				if _, ok := valMap[v.Address]; !ok {
					validator.InactiveValidator(v.Address, int64(types.Inactive), 0, 0, time.Time{})
				}
			}
		}
	}

	return nil
}

func CreateConsensusState(chainID string, cs *tmctypes.ResultConsensusState) error {
	return block.UpdateConsensusState(chainID, cs)
}

// SyncLock 同步时锁定，同一个时间只会有一个同步协程
func SyncLock(key string) bool {
	key = "lock_" + key
	value := "1"

	qs, err := model.QmoonStatusByKey(db.Db, utils.NullString(key))
	if err != nil {
		if err == sql.ErrNoRows {
			qs = &model.QmoonStatus{
				Key:   utils.NullString(key),
				Value: utils.NullString(value),
			}
			err := qs.Insert(db.Db)
			return err == nil
		}
	}

	if qs.Value.String == "1" {
		return false
	}

	s := fmt.Sprintf("update public.qmoon_status set value='%s' where key='%s' and value='0'", value, key)
	log.Printf(s)
	r, err := db.Db.Exec(s)
	if err != nil {
		return false
	}

	num, err := r.RowsAffected()
	if err != nil {
		return false
	}

	return num == 1
}

func SyncUnlock(key string) bool {
	key = "lock_" + key

	value := "0"
	s := fmt.Sprintf("update public.qmoon_status set value='%s' where key='%s' and value='1'", value, key)
	log.Printf(s)
	r, err := db.Db.Exec(s)
	if err != nil {
		return false
	}

	num, err := r.RowsAffected()
	if err != nil {
		return false
	}

	return num == 1
}
