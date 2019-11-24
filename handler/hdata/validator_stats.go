// Copyright 2018 The QOS Authors

package hdata

import (
	"github.com/QOSGroup/qmoon/models"
	"net/http"

	"github.com/QOSGroup/qmoon/handler/middleware"
	"github.com/QOSGroup/qmoon/lib"
	"github.com/QOSGroup/qmoon/types"
	"github.com/gin-gonic/gin"
)

const validatorVotingPowerUrl = NodeProxy + "/validators/:address/votingPower"
const validatorStatUrl = NodeProxy + "/validators/:address/stats"

func init() {
	hdataHander[validatorVotingPowerUrl] = ValidatorVotingPowerGinRegister
	hdataHander[validatorStatUrl] = validatorStatGinRegister
}

func validatorStatGinRegister(r *gin.Engine) {
	r.GET(validatorStatUrl, middleware.ApiAuthGin(), validatorStatGin())
}
func ValidatorVotingPowerGinRegister(r *gin.Engine) {
	r.GET(validatorVotingPowerUrl, middleware.ApiAuthGin(), validatorVotingPowerGin())
}

func validatorStatGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		result := types.ResultValidatorStats{
			Operations: make([]*types.ValidatorOperations, 0),
			Proposed: make([]int64, 0),
			Missed: make([]int64, 0),
			Evidence:make([]int64, 0),
		}
		node, err := GetNodeFromUrl(c)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCMethodNotFoundError(""))
			return
		}

		val, err := models.ValidatorByStakeAddress(node.ChainID, c.Param("address"));
		if err != nil {
			c.JSON(http.StatusOK, types.RPCServerError("", err))
			return
		}
		prop, err := models.BlocksByProposer(node.ChainID,val.Address)
		if err == nil && prop != nil && len(prop) > 0 {
			for _, p := range prop {
				result.Proposed = append(result.Proposed, p.Height)
			}
		}
		missings, err := node.Missings(val.Address)
		if err == nil && missings != nil && len(missings) > 0 {
			for _, m := range missings {
				result.Missed = append(result.Missed, m.Height)
			}
		}

		result.Operations, _ = models.TxByAddressAndType(node.ChainID, val.Owner, 0, 0, 0, 100000,
			"stake/txs/TxCreateValidator", "stake/txs/TxModifyValidator", "stake/txs/TxRevokeValidator", "stake/txs/TxActiveValidator")

		c.JSON(http.StatusOK, types.NewRPCSuccessResponse(lib.Cdc, "", result))
	}
}

func validatorVotingPowerGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		node, err := GetNodeFromUrl(c)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCMethodNotFoundError(""))
			return
		}

		val, err := models.ValidatorByStakeAddress(node.ChainID, c.Param("address"));
		if err != nil {
			c.JSON(http.StatusOK, types.RPCServerError("", err))
			return
		}
		res, err := models.QueryValidatorVotingPower(node.ChainID, val.Address, 1000)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCServerError("", err))
			return
		}

		c.JSON(http.StatusOK, types.NewRPCSuccessResponse(lib.Cdc, "", res))
	}
}
