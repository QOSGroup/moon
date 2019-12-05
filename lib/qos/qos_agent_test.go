package qos

import (
	"github.com/gin-gonic/gin/json"
	"log"
	"testing"
)

func TestTx(t *testing.T) {
	proposals, err := NewQosCli("").QueryTx("47.103.79.28:26657", "C4CE4728774B67C064DF3274063A9035C95D0E92AD3073C717047B7A48486CBD")
	bytes, err := json.Marshal(proposals)
	log.Printf("res:%+v, err:%+v", string(bytes), err)
}

func TestQueryProposals(t *testing.T) {
	proposals, err := NewQosCli("").QueryProposals("47.103.79.28:26657")
	//fmt.Println(proposals[0])
	bytes, err := json.Marshal(proposals)
	log.Printf("res:%+v, err:%+v", string(bytes), err)
}

func TestQueryProposal(t *testing.T) {
	proposals, err := NewQosCli("").QueryProposal("39.97.234.227:26657", 1)
	bytes, err := json.Marshal(proposals)
	log.Printf("res:%+v, err:%+v", string(bytes), err)
}

func TestQueryVotes(t *testing.T) {
	proposals, err := NewQosCli("").QueryVotes("39.97.234.227:26657", 1)
	bytes, err := json.Marshal(proposals)
	log.Printf("res:%+v, err:%+v", string(bytes), err)
}

func TestQueryDeposits(t *testing.T) {
	proposals, err := NewQosCli("").QueryDeposits("39.97.234.227:26657", 1)
	bytes, err := json.Marshal(proposals)
	log.Printf("res:%+v, err:%+v", string(bytes), err)
}

func TestQueryTally(t *testing.T) {
	proposals, err := NewQosCli("").QueryTally("39.97.234.227:26657", 1)
	bytes, err := json.Marshal(proposals)
	log.Printf("res:%+v, err:%+v", string(bytes), err)
}

func TestQueryDelegationsWithValidator(t *testing.T) {
	proposals, err := NewQosCli("").QueryDelegationsWithValidator("47.103.79.28:26657", "qosval19hrl38w5lm6sklw2hzrzrjtsxudpy8hyfaea3e")
	bytes, err := json.Marshal(proposals)
	log.Printf("res:%+v, err:%+v", string(bytes), err)
}
