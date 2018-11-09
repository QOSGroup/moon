// Copyright 2018 The QOS Authors

package hdata

import (
	"log"
	"net/http"

	"github.com/QOSGroup/qmoon/handler/middleware"
	"github.com/QOSGroup/qmoon/lib"
	"github.com/QOSGroup/qmoon/service/tx"
	"github.com/QOSGroup/qmoon/types"
	"github.com/gin-gonic/gin"
)

const sequenceQueryUrl = "/sequence"

func init() {
	hdataHander[sequenceQueryUrl] = SequenceQueryGinRegister
}

// SequenceQueryGinRegister 注册sequenceQuery
func SequenceQueryGinRegister(r *gin.Engine) {
	r.GET(nodeProxy+sequenceQueryUrl, middleware.ApiAuthGin(), sequenceQueryGin())
}

func sequenceQueryGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		nt, err := getNodeFromUrl(c)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCMethodNotFoundError(""))
			return
		}

		gs, err := nt.AppState(lib.Cdc)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCInternalError("", err))
			return
		}

		ctx := tx.NewClient(nt.BaseURL)
		var result types.ResultSequence
		for _, v := range gs.QCPs {
			log.Printf("sequenceQuery url:%s, chainID:%s", nt.BaseURL, v.Name)
			in, _ := ctx.SequenctIn(v.Name)
			out, _ := ctx.SequenctOut(v.Name)
			result.Apps = append(result.Apps, &types.Sequence{Name: v.Name, In: in, Out: out})
		}

		c.JSON(http.StatusOK, types.NewRPCSuccessResponse(lib.Cdc, "", result))
	}
}
