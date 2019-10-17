// Copyright 2018 The QOS Authors

package service

import (
	"time"

	"github.com/QOSGroup/qmoon/lib/cache"
	"github.com/QOSGroup/qmoon/types"
)

const chainStatusCache = "ChainStatusCache"

func (n Node) ChainStatus(cached bool) (*types.ResultStatus, error) {
	result := &types.ResultStatus{}
	if cached {
		d, ok := cache.Get(chainStatusCache)
		if ok {
			if v, okk := d.(*types.ResultStatus); okk {
				return v, nil
			}
		}
	}

	cs, err1 := n.ConsensusState()
	if err1 != nil {
		result.ConsensusState = &types.ResultConsensusState{}
	} else {
		result.ConsensusState = cs
	}

	vs, err2 := n.Validators()
	if err2 == nil {
		result.TotalValidators = int64(len(vs))
	}

	d, err := n.BlockTimeAvg(100)
	if err == nil {
		result.BlockTimeAvg = d.String()
	}

	lb, err3 := n.LatestBlock()
	if err3 == nil {
		result.TotalTxs = lb.TotalTxs
	}
	result.ConsensusState.ChainID = n.ChainID

	if err3 == nil && err2 == nil {
		cache.Set(chainStatusCache, result, time.Second*1)
	}

	return result, nil
}
