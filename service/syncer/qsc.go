// Copyright 2018 The QSC Authors

package syncer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/QOSGroup/qmoon/lib"
	"github.com/QOSGroup/qmoon/models"
	"github.com/QOSGroup/qmoon/service"
	"github.com/QOSGroup/qmoon/service/metric"
	"github.com/QOSGroup/qmoon/types"
	"github.com/QOSGroup/qmoon/utils"
	tmtypes "github.com/tendermint/tendermint/types"
)

type QSC struct {
	node  *service.Node
	tmcli *lib.TmClient
}

func (s QSC) RpcPeers(ctx context.Context) error {
	netinfo, err := s.tmcli.NetInfo()
	if err != nil {
		return err
	}

	//var peers []string
	for _, peer := range netinfo.Peers {
		remoteIp := peer.RemoteIP
		u, err := url.Parse(peer.NodeInfo.Other.RPCAddress)
		if err != nil {
			log.Printf("[QOS] Parse Peer RPCAddress err:%s", err.Error())
			continue
		}

		u.Host = fmt.Sprintf("%s:%s", remoteIp, u.Port())
		//peers = append(peers, u.String())

		_ = models.CreateNetworkOrUpdate(s.node.ChanID, &models.Network{Remote: u.String()})
	}

	return nil
}

// BlockLoop 同步块
func (s QSC) BlockLoop(ctx context.Context) error {
	key := "lock_" + s.node.ChanID + "-" + LockTypeBlock

	if !Lock(key) {
		log.Printf("[Sync] BlockLoop %v err, has been locked.", s.node.ChanID)
		return nil
	}
	defer Unlock(key)

	var height int64 = 1
	latest, err := s.node.LatestBlock()
	if err == nil && latest != nil {
		height = latest.Height + 1
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			b, err := s.tmcli.RetrieveBlock(&height)
			if err != nil {
				time.Sleep(time.Millisecond * 1000)
				continue
			}

			if err := s.block(b); err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			s.Validator(height, b.Header.Time)

			height += 1
		}
	}

	return nil
}

// block
func (s QSC) block(b *types.Block) error {
	err := s.node.CreateBlock(b)
	if err != nil {
		return err
	}

	err = s.tx(b)
	// TODO delete block

	err = s.node.SaveBlockValidator(b.Precommits)
	// TODO delete block and tx

	return nil
}

func (s QSC) tx(b *types.Block) error {
	if len(b.Txs) == 0 {
		return nil
	}

	var txs []string
	for _, v := range b.Txs {
		txs = append(txs, utils.Base64En(v))
	}
	txres, err := lib.NewQstarsAgentCli("").Txs(txs)
	if err != nil {
		return err
	}

	for k, v := range txres {
		hash := tmtypes.Tx(b.Txs[k]).Hash()
		mt := &models.Tx{}
		mt.Height = b.Header.Height
		mt.Index = int64(k)
		mt.Hash = strings.ToUpper(hex.EncodeToString(hash))
		mt.TxType = v.Tx.Type
		mt.OriginTx = txs[k]
		if d, err := json.Marshal(v.Tx); err == nil {
			mt.JsonTx = string(d)
		}
		mt.Time = b.Header.Time

		txResult, err := s.tmcli.RetrieveTx(hash)
		if err == nil {
			mt.TxStatus = int(txResult.TxStatus)
			mt.GasWanted = txResult.GasWanted
			mt.GasUsed = txResult.GasUsed
			mt.Log = txResult.Log
		}
		if err := mt.Insert(s.node.ChanID); err != nil {
			log.Printf("tx insert data:%+v, err:%v", mt, err.Error())
			return err
		}
	}

	return nil
}

func (s QSC) Validator(height int64, t time.Time) error {
	vals, err := s.tmcli.Validator(height)
	if err != nil {
		time.Sleep(time.Millisecond * 100)
		return err
	}

	for _, val := range vals {
		_ = s.node.CreateValidator(val)
	}

	valMap := make(map[string]types.Validator)
	for _, v := range vals {
		valMap[v.Address] = v
	}

	allVals, err := s.node.Validators()
	if err == nil {
		for _, v := range allVals {
			if v.Status == types.Active {
				if _, ok := valMap[v.Address]; !ok {
					_ = s.node.InactiveValidator(v.Address, int(types.Inactive), height, time.Time{})
				}
			} else {
				if _, ok := valMap[v.Address]; ok {
					_ = s.node.InactiveValidator(v.Address, int(types.Active), height, time.Time{})
				}
			}
		}
	}
	metric.ValidatorVotingPower(s.node.ChanID, t, vals)

	return nil
}

func (s QSC) ConsensusStateLoop(ctx context.Context) error {
	key := "lock_" + s.node.ChanID + "-" + LockTypeConsensusState

	if !Lock(key) {
		log.Printf("[Sync] ConsensusStateLoop %v err, has been locked.", s.node.ChanID)

		return nil
	}
	defer Unlock(key)

	for {
		time.Sleep(SyncConsensusStateDuration)
		select {
		case <-ctx.Done():
			return nil
		default:
			cs, err := s.tmcli.ConsensusState()
			if err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			if err := s.node.UpdateConsensusState(cs); err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}
		}
	}

	return nil
}

func (s QSC) PeerLoop(ctx context.Context) error {
	key := "lock_" + s.node.ChanID + "-" + LockTypePeer

	if !Lock(key) {
		log.Printf("[Sync] PeerLoop %v err, has been locked.", s.node.ChanID)
		return nil
	}
	defer Unlock(key)

	for {
		time.Sleep(SyncPeerDuration)
		select {
		case <-ctx.Done():
			return nil
		default:
			b, err := s.tmcli.NetInfo()
			if err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			if b != nil {
				if err := s.node.CreatePeers(b.Peers); err != nil {
					time.Sleep(time.Millisecond * 100)
					continue
				}
			}
		}
	}

	return nil
}
