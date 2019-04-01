// Copyright 2018 The QOS Authors

// Package pkg comments for pkg block
// service 块相关数据封装
package service

import (
	"errors"
	"time"

	"github.com/QOSGroup/qmoon/models"
	"github.com/QOSGroup/qmoon/types"
)

func convertToBlock(mb *models.Block) *types.ResultBlockBase {
	return &types.ResultBlockBase{
		ID:             mb.Id,
		ChainID:        mb.ChainId,
		Height:         mb.Height,
		NumTxs:         mb.NumTxs,
		TotalTxs:       mb.TotalTxs,
		Time:           types.ResultTime(mb.Time),
		DataHash:       mb.DataHash,
		ValidatorsHash: mb.ValidatorsHash,
	}
}

// Latest最新的块
func (n Node) LatestBlock() (*types.ResultBlockBase, error) {
	mbs, err := models.Blocks(n.ChanID, &models.BlockOption{Offset: 0, Limit: 1})
	if err != nil {
		return nil, err
	}

	if len(mbs) == 0 {
		return nil, errors.New("not found")
	}

	return convertToBlock(mbs[0]), nil
}

// Retrieve 块查询
func (n Node) RetrieveBlock(height int64) (*types.ResultBlockBase, error) {
	mbs, err := models.Blocks(n.ChanID, &models.BlockOption{Height: height, Offset: 0, Limit: 1})
	if err != nil {
		return nil, err
	}

	if len(mbs) != 1 {
		return nil, errors.New("not found")
	}

	return convertToBlock(mbs[0]), err
}

// Search 块查询
func (n Node) Blocks(minHeight, maxHeight int64) ([]*types.ResultBlockBase, error) {
	mbs, err := models.Blocks(n.ChanID, &models.BlockOption{MinHeight: minHeight, MaxHeight: maxHeight})
	if err != nil {
		return nil, err
	}

	var res []*types.ResultBlockBase
	for _, v := range mbs {
		res = append(res, convertToBlock(v))
	}

	return res, err
}

// Search 最近N块平均打快时间
func (n Node) BlockTimeAvg(blockNum int) (time.Duration, error) {
	mbs, err := models.Blocks(n.ChanID, &models.BlockOption{Offset: 0, Limit: blockNum})
	if err != nil {
		return 0, err
	}

	if len(mbs) <= 1 {
		return 0, nil
	}

	var duration int64
	num := int64(0)
	for k := 1; k < len(mbs); k++ {
		duration += int64(mbs[k-1].Time.Sub(mbs[k].Time))
		num++
	}

	return time.Duration(duration / num), err
}

// HasTx 有交易的块
func (n Node) HasTxBlocks(minHeight, maxHeight int64) ([]*types.ResultBlockBase, error) {
	mbs, err := models.Blocks(n.ChanID, &models.BlockOption{
		MinHeight: minHeight, MaxHeight: maxHeight, NumTxs: 1})
	if err != nil {
		return nil, err
	}

	var res []*types.ResultBlockBase
	for _, v := range mbs {
		res = append(res, convertToBlock(v))
	}

	return res, err
}

func (n Node) CreateBlock(b *types.Block) error {
	block := &models.Block{}
	block.Height = b.Header.Height
	block.NumTxs = b.Header.NumTxs
	block.TotalTxs = b.Header.TotalTxs
	block.Time = b.Header.Time
	block.DataHash = b.Header.DataHash
	block.ValidatorsHash = b.Header.ValidatorsHash
	block.ProposerAddress = b.Header.ProposerAddress
	if err := block.Insert(n.ChanID); err != nil {
		return err
	}

	return nil
}

func (n Node) CreateEvidence(b *types.Block) error {
	if b.EvidenceList.Evidences == nil || len(b.EvidenceList.Evidences) == 0 {
		return nil
	}

	return models.CreateEvidences(n.ChanID, b.EvidenceList)
}
