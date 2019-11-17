// Copyright 2018 The QOS Authors

package service

import (
	"github.com/QOSGroup/qmoon/lib/cache"
	"github.com/QOSGroup/qmoon/lib/qos"
	"github.com/QOSGroup/qmoon/types"
	"time"
)

const chainStatusCache = "ChainStatusCache"
const LatestHeightKey = "community_fee_pool_key"

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


	status, err := qos.NewQosCli("").QueryStatus(n.BaseURL)
	if err != nil {
		return nil, err
	}

	if status!=nil {
		result.Height = status.SyncInfo.LatestBlockHeight
		cache.Set(LatestHeightKey, result.Height,  time.Second*7)
		blc, err := n.BlockByHeight(result.Height)
		if err == nil {
			result.Block = blc
		}

		lb, err3 := n.BlockByHeight(result.Height)
		if err3 == nil {
			result.TotalTxs = lb.TotalTxs
			result.Proposer = lb.Proposer
			result.Votes = lb.Votes
		}

		vs, err2 := n.Validators(result.Height)
		if err2 == nil {
			result.TotalValidators = int64(len(vs))
		}

	}

	cs, err1 := n.ConsensusState()
	if err1 != nil {
		result.ConsensusState = &types.ResultConsensusState{}
	} else {
		result.ConsensusState = cs
		// latestHeight,_ = strconv.ParseInt(cs.Height, 10, 64)
	}



	d, err := n.BlockTimeAvg(100)
	if err == nil {
		result.BlockTimeAvg = d.String()
	}

	// lb, err3 := n.LatestBlock()

	result.ConsensusState.ChainID = n.ChainID

	cache.Set(chainStatusCache, result, time.Second*1)

	return result, nil
}
