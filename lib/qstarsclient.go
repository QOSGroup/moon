// Copyright 2018 The QOS Authors

package lib

import (
	"os"

	"github.com/QOSGroup/qbase/account"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qmoon/types"
	qstarscontext "github.com/QOSGroup/qstars/client/context"
	"github.com/tendermint/go-amino"
)

var qstarscs map[string]*QstarsClient

func init() {
	qstarscs = make(map[string]*QstarsClient)
}

type QstarsClient struct {
	c   qstarscontext.CLIContext
	cdc *amino.Codec
}

func NewQstarsClient(nodeURI string) *QstarsClient {

	qc, ok := qstarscs[nodeURI]
	if ok {
		return qc
	}
	qc = &QstarsClient{}
	cdc := MakeCodec()
	qc.c = qstarscontext.NewCLIContext1(nodeURI).
		WithCodec(cdc).
		WithLogger(os.Stdout)
	qc.cdc = cdc
	qstarscs[nodeURI] = qc

	return qc
}

func (c QstarsClient) MarshalJSONIndent(o interface{}, prefix, indent string) ([]byte, error) {
	return c.cdc.MarshalJSONIndent(o, prefix, indent)
}

func (c QstarsClient) UnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	return c.cdc.UnmarshalBinaryBare(bz, ptr)
}

func (c QstarsClient) UnmarshalJSON(bz []byte, ptr interface{}) error {
	return c.cdc.UnmarshalJSON(bz, ptr)
}

func (c QstarsClient) MarshalJSON(o interface{}) ([]byte, error) {
	return c.cdc.MarshalJSON(o)
}

// QuerySequenceIn query sequence in
func (c QstarsClient) SequenctIn(chainID string) (int64, error) {
	return c.c.QuerySequenceIn(chainID)
}

// QuerySequenceOut query sequence out
func (c QstarsClient) SequenctOut(chainID string) (int64, error) {
	return c.c.QuerySequenceOut(chainID)
}

// QueryAccount query account by addr
func (c QstarsClient) QueryAccount(addr string) (*types.ResultAccount, error) {
	address, err := qbasetypes.GetAddrFromBech32(addr)
	if err != nil {
		return nil, err
	}
	key := account.AddressStoreKey(address)

	acc, err := c.c.GetAccount(key, c.cdc)
	if err != nil {
		return nil, err
	}

	result := &types.ResultAccount{}
	result.Address = acc.AccountAddress.String()
	result.Nonce = acc.Nonce

	result.Coins = append(result.Coins, &qbasetypes.BaseCoin{Name: "qos", Amount: acc.QOS})
	for _, v := range acc.QSCs {
		result.Coins = append(result.Coins, v)
	}

	return result, nil
}
