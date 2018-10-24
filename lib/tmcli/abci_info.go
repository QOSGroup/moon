package tmcli

import (
	"context"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func init() {

}

const abciInfoURI = "abci_info"

type AbciInfoService service

func (s *AbciInfoService) Retrieve(ctx context.Context) (*tmctypes.ResultABCIInfo, error) {
	u := abciInfoURI

	u, err := addOptions(u, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var res tmctypes.ResultABCIInfo
	_, err = s.client.Do(ctx, req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
