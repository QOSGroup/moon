// Copyright 2018 The QOS Authors

package service

import (
	"time"

	"github.com/QOSGroup/qmoon/lib/cache"
	"github.com/QOSGroup/qmoon/service/block"
	"github.com/QOSGroup/qmoon/service/validator"
	"github.com/QOSGroup/qmoon/types"
)

const chainStatusCache = "ChainStatusCache"

func ChainStatus(chainID string, cached bool) (*types.ResultStatus, error) {
	result := &types.ResultStatus{}

	if cached {
		d, ok := cache.Get(chainStatusCache)
		if ok {
			if v, okk := d.(*types.ResultStatus); okk {
				return v, nil
			}
		}
	}

	cs, err1 := block.RetrieveConsensusState(chainID)
	result.ConsensusState = cs

	vs, err2 := validator.ListValidatorByChain(chainID)
	result.TotalValidators = int64(len(vs))

	if err1 == nil && err2 == nil {
		cache.Set(chainStatusCache, result, time.Second*1)
	}

	return result, nil
}
