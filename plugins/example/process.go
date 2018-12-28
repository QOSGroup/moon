// Copyright 2018 The QOS Authors

package example

import (
	"database/sql"

	qbasetxs "github.com/QOSGroup/qbase/txs"
	"github.com/gin-gonic/gin"
	tmtypes "github.com/tendermint/tendermint/types"
)

type EgPlugin struct{}

func (ttp EgPlugin) DbInit(driveName string, db *sql.DB) error {

	return DbInit(driveName, db)
}

func (ttp EgPlugin) DbClear(driveName string, db *sql.DB) error {
	return DbClear(driveName, db)
}

func (ttp EgPlugin) Type() string {
	return "Eg"
}

func (ttp EgPlugin) RegisterGin(r *gin.Engine) {
	ExampleGinRegister(r)
}

func (ttp EgPlugin) Parse(blockHeader tmtypes.Header, itx qbasetxs.ITx) (typeName string, hit bool, err error) {
	return "", false, nil
}

func (ttp EgPlugin) Doctor() error {
	return nil
}
