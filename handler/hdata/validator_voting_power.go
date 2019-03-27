// Copyright 2018 The QOS Authors

package hdata

import (
	"net/http"
	"strconv"
	"time"

	"github.com/QOSGroup/qmoon/handler/middleware"
	"github.com/QOSGroup/qmoon/lib"
	"github.com/QOSGroup/qmoon/service/metric"
	"github.com/QOSGroup/qmoon/types"
	"github.com/QOSGroup/qmoon/utils"
	"github.com/gin-gonic/gin"
)

const validatorVotingPowerUrl = NodeProxy + "/validators/:address/votingPower"

func init() {
	hdataHander[validatorVotingPowerUrl] = ValidatorVotingPowerGinRegister
}

// ValidatorVotingPowerGinRegister 注册
func ValidatorVotingPowerGinRegister(r *gin.Engine) {
	r.GET(validatorVotingPowerUrl, middleware.ApiAuthGin(), validatorVotingPowerGin())
}

func validatorVotingPowerGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		node, err := GetNodeFromUrl(c)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCMethodNotFoundError(""))
			return
		}

		address := lib.Bech32AddressToHex(c.Param("address"))
		now := time.Now()
		start := utils.NDaysAgo(now, 28).Unix()
		end := now.Unix()
		step := int64(60 * 60 * 24)

		if d, err := strconv.ParseInt(c.Query("start"), 10, 64); err == nil {
			start = d
		}
		if d, err := strconv.ParseInt(c.Query("end"), 10, 64); err == nil {
			end = d
		}
		if d, err := strconv.ParseInt(c.Query("step"), 10, 64); err == nil {
			step = d
		}
		res, err := metric.QueryValidatorVotingPower(node.ChanID, address,
			time.Unix(start, 0), time.Unix(end, 0), time.Second*time.Duration(step))
		if err != nil {
			c.JSON(http.StatusOK, types.RPCServerError("", err))
			return
		}

		c.JSON(http.StatusOK, types.NewRPCSuccessResponse(lib.Cdc, "", res))
	}
}
