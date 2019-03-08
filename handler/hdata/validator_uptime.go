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
	"github.com/gin-gonic/gin"
)

const validatorUptimeUrl = NodeProxy + "/validators/:address/uptime"

func init() {
	hdataHander[validatorUptimeUrl] = ValidatorUptimeGinRegister
}

// ValidatorVotingPowerGinRegister 注册
func ValidatorUptimeGinRegister(r *gin.Engine) {
	r.GET(validatorUptimeUrl, middleware.ApiAuthGin(), validatorUptimeGin())
}

func validatorUptimeGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		node, err := GetNodeFromUrl(c)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCMethodNotFoundError(""))
			return
		}

		address := lib.Bech32AddressToHex(c.Param("address"))
		start, err := strconv.ParseInt(c.Query("start"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCInvalidParamsError("", err))
			return
		}
		end, err := strconv.ParseInt(c.Query("end"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCInvalidParamsError("", err))
			return
		}
		step, err := strconv.ParseInt(c.Query("step"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCInvalidParamsError("", err))
			return
		}

		res, err := metric.QueryValidatorUptime(node.ChanID, address,
			time.Unix(start, 0), time.Unix(end, 0), time.Second*time.Duration(step))
		if err != nil {
			c.JSON(http.StatusOK, types.RPCServerError("", err))
			return
		}

		c.JSON(http.StatusOK, types.NewRPCSuccessResponse(lib.Cdc, "", res))
	}
}
