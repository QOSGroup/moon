package models

import (
	"github.com/QOSGroup/qmoon/models/errors"
)

type ConsensusState struct {
	Id              int64  `xorm:"pk autoincr BIGINT"`
	ChainId         string `xorm:"-"`
	Height          string `xorm:"TEXT"`
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

func RetrieveVotesByHeight(chainID string, height int64) (string, error) {
	x, err := GetNodeEngine(chainID)
	if err != nil {
		return "Not available", err
	}
	cs := &ConsensusState{Height:string(height)}
	has, err := x.Get(cs)
	if err != nil {
		return "Not available", err
	}

	if !has {
		return "Not available", errors.NotExist{Obj: "ConsensusStateLoop"}
	}

	return string(cs.PrevotesNum) + "/" + string(cs.PrecommitsNum), nil
}
