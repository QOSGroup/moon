// Copyright 2018 The QOS Authors

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/QOSGroup/qmoon/config"
	"github.com/QOSGroup/qmoon/db/migrations"
	"github.com/QOSGroup/qmoon/db/model"
	"github.com/QOSGroup/qmoon/utils"

	_ "github.com/lib/pq"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func checkDriverSupport(d string) bool {
	if d == "mysql" {
		return true
	}

	if d == "postgres" {
		return true
	}

	return false
}

func InitDb(cfg *config.DBConfig, logger *log.Logger) error {
	if !checkDriverSupport(cfg.DriverName) {
		return errors.New("unsupported db driver")
	}

	d, err := sql.Open(cfg.DriverName, cfg.DataSource())
	if err != nil {
		return err
	}

	if err := d.Ping(); err != nil {
		return err
	}

	Db = d

	return nil
}

// checkDatabase 检查数据库是否存在
func checkDatabase(dbName string, db *sql.DB) (bool, error) {
	s := fmt.Sprintf("select count(*) from pg_catalog.pg_database where datname='%s';", dbName)
	var count int
	err := db.QueryRow(s).Scan(&count)

	if err != nil {
		return false, err
	}

	return count == 1, nil
}

// createDatabase 创建数据库
func createDatabase(dbName string, db *sql.DB) error {
	s := fmt.Sprintf("create database %s ENCODING 'UTF8' TEMPLATE template0;", dbName)

	_, err := db.Query(s)

	return err
}

// DropDatabase 删除数据库
func dropDatabase(dbName string, db *sql.DB) error {
	var err error

	closeConn := fmt.Sprintf("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity "+
		"WHERE datname='%s' AND pid<>pg_backend_pid();", dbName)
	_, err = db.Query(closeConn)
	if err != nil {
		return err
	}

	s := fmt.Sprintf("drop database %s;", dbName)
	_, err = db.Query(s)

	return err
}

type testDb struct {
	db          *sql.DB
	tmpDatabase string
}

func (dt *testDb) Close() {
	defer dt.db.Close()
	dropDatabase(dt.tmpDatabase, dt.db)
}

func NewTestDb(m *testing.M) *testDb {
	res := &testDb{}

	dbCfg := config.TestDBConfig()
	dbCfg.Database = ""

	err := InitDb(dbCfg, log.New(os.Stderr, "", log.LstdFlags))
	if err != nil {
		panic(err)
	}
	res.db = Db

	dbCfg.Database = fmt.Sprintf("qmoon_test_%d", time.Now().Nanosecond())
	ok, err := checkDatabase(dbCfg.Database, Db)
	if err != nil {
		panic(err)
	}
	if !ok {
		err := createDatabase(dbCfg.Database, res.db)
		if err != nil {
			panic(err)
		}

		err = InitDb(dbCfg, log.New(os.Stderr, "", log.LstdFlags))
		if err != nil {
			panic(err)
		}
	}

	err = migrations.Up(dbCfg.DriverName, Db)
	if err != nil {
		panic(err)
	}

	model.XOLog = utils.SQLLog(os.Stdout)

	res.tmpDatabase = dbCfg.Database

	return res
}
