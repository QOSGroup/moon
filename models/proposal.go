package models

import (
	"time"

	"github.com/go-xorm/xorm"
)

type Proposal struct {
	Id                  int64     `xorm:"pk autoincr BIGINT"`
	ProposalID          int64     `xorm:"unique(proposal_proposal_idx) BIGINT"`
	Title               string    `xorm:"TEXT"`
	Description         string    `xorm:"TEXT"`
	Type                string    `xorm:"TEXT"`
	Status              string    `xorm:"TEXT"`
	SubmitTime          time.Time `xorm:"-"`
	VotingStartTime     time.Time `xorm:"-"`
	VotingEndTime       time.Time `xorm:"-"`
	SubmitTimeUnix      int64
	VotingStartTimeUnix int64
	VotingEndTimeUnix   int64
	TotalDeposit        int64  `xorm:"BIGINT"`
	ChainID             string `xorm:"-"`
	// UpdatedAt           time.Time `xorm:"-"`
	// UpdatedAtUnix       int64
	// CreatedAt           time.Time `xorm:"-"`
	// CreatedAtUnix       int64
}

func (n *Proposal) BeforeInsert() {
	n.SubmitTimeUnix = n.SubmitTime.Unix()
	n.VotingStartTimeUnix = n.VotingStartTime.Unix()
	n.VotingEndTimeUnix = n.VotingEndTime.Unix()
}

func (n *Proposal) BeforeUpdate() {
	n.SubmitTimeUnix = n.SubmitTime.Unix()
	n.VotingStartTimeUnix = n.VotingStartTime.Unix()
	n.VotingEndTimeUnix = n.VotingEndTime.Unix()
}

func (n *Proposal) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "submit_at_unix":
		n.SubmitTime = time.Unix(n.SubmitTimeUnix, 0).Local()
	case "voting_start_at_unix":
		n.VotingStartTime = time.Unix(n.VotingStartTimeUnix, 0).Local()
	case "voting_end_at_unix":
		n.VotingEndTime = time.Unix(n.VotingEndTimeUnix, 0).Local()

	}
}

func (n *Proposal) Insert(chainID string) error {
	x, err := GetNodeEngine(chainID)
	if err != nil {
		return err
	}

	_, err = x.Insert(n)
	if err != nil {
		return err
	}

	return nil
}

func (pro *Proposal) Update(chainID string) error {
	x, err := GetNodeEngine(chainID)
	if err != nil {
		return err
	}

	_, err = x.ID(pro.Id).Update(pro)
	if err != nil {
		return err
	}

	return nil
}

func (pro *Proposal) UpdateStatus(chainID string) error {
	x, err := GetNodeEngine(chainID)
	if err != nil {
		return err
	}

	_, err = x.ID(pro.Id).
		Cols("status").
		Update(pro)
	if err != nil {
		return err
	}

	return nil
}

func Proposals(chainID string) ([]*Proposal, error) {
	x, err := GetNodeEngine(chainID)
	if err != nil {
		return nil, err
	}
	var pros = make([]*Proposal, 0)
	return pros, x.Find(&pros)
}
