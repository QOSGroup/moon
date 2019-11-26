package models

import (
	"github.com/QOSGroup/qmoon/models/errors"
	"strconv"
)

type ConsensusState struct {
	Id              int64  `xorm:"pk autoincr BIGINT"`
	ChainId         string `xorm:"-"`
	Height          string `xorm:"unique(block_height) TEXT"`
	Round           string `xorm:"TEXT"`
	Step            string `xorm:"TEXT"`
	PrevotesNum     int64  `xorm:"BIGINT"`
	PrevotesValue   string `xorm:"TEXT"`
	PrecommitsNum   int64  `xorm:"BIGINT"`
	PrecommitsValue string `xorm:"TEXT"`
	StartTime       string `xorm:"TEXT"`
}

func (cs *ConsensusState) Insert(chainID string) error {
	x, err := GetNodeEngine(chainID)
	if err != nil {
		return err
	}

	_, err = x.Insert(cs)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ConsensusState) Update(chainID string) error {
	x, err := GetNodeEngine(chainID)
	if err != nil {
		return err
	}

	_, err = x.ID(cs.Id).Update(cs)
	if err != nil {
		return err
	}

	return nil
}

func RetrieveConsensusState(chainID string) (*ConsensusState, error) {
	x, err := GetNodeEngine(chainID)
	if err != nil {
		return nil, err
	}
	cs := &ConsensusState{}
	has, err := x.Get(cs)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.NotExist{Obj: "ConsensusStateLoop"}
	}

	return cs, nil
}

func RetrieveConsensusStateByHeight(chainID string, height string) (*ConsensusState, error) {
	x, err := GetNodeEngine(chainID)
	if err != nil {
		return nil, err
	}
	cs := &ConsensusState{Height:height}
	has, err := x.Get(cs)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.NotExist{Obj: "ConsensusStateLoop"}
	}

	return cs, nil
}

func RetrieveVotesByHeight(chainID string, height int64) (string, error) {
	cs, err := RetrieveConsensusStateByHeight(chainID, strconv.FormatInt(height, 10))
	if err != nil {
		return "Not available", err
	}

	return strconv.FormatInt(cs.PrevotesNum, 10) + "/" + strconv.FormatInt(cs.PrecommitsNum, 10), nil
}
