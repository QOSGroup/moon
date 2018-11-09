// Copyright 2018 The QOS Authors

package worker

import (
	"sync"

	"github.com/QOSGroup/qmoon/service"
	"github.com/tendermint/tendermint/rpc/client"
)

// SyncPeersLoop 同步peer节点信息
func SyncPeersLoop() {
	wg := sync.WaitGroup{}

	nodes, err := service.AllNodes()
	if err != nil {
		return
	}

	for _, v := range nodes {
		//log.Printf("--SyncBlockLoop start chanID:%s", chanID)
		wg.Add(1)
		go func() {
			defer wg.Done()
			SyncPeer(v.ChanID, v.BaseURL)
		}()
	}

	wg.Wait()
}

// SyncPeer 同步p2p中peer信息
func SyncPeer(chanID, remote string) error {
	tmc := client.NewHTTP(remote, "/websocket")

	b, err := tmc.NetInfo()
	if err != nil {
		return err
	}

	if b == nil {
		return nil
	}

	return service.SavePeers(chanID, b.Peers)
}
