// Copyright 2018 The QOS Authors

package hdata

import (
	"net/http"

	"github.com/QOSGroup/qmoon/handler/middleware"
	"github.com/QOSGroup/qmoon/lib"
	"github.com/QOSGroup/qmoon/service"
	"github.com/QOSGroup/qmoon/types"
	"github.com/gin-gonic/gin"
)

const peersUrl = "/peers"

func init() {
	hdataHander[peersUrl] = PeersGinRegister
}

// PeersGinRegister 注册peers
func PeersGinRegister(r *gin.Engine) {
	r.GET(nodeProxy+peersUrl, middleware.ApiAuthGin(), peersGin())
}

func peersGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		nt, err := getNodeFromUrl(c)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCMethodNotFoundError(""))
			return
		}

		result, err := service.ListPeers(nt.ChanID)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCServerError("", err))
			return
		}

		c.JSON(http.StatusOK, types.NewRPCSuccessResponse(lib.Cdc, "", result))
	}
}
