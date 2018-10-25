// Copyright 2018 The QOS Authors

// Package pkg comments for pkg admin
// admin 管理后台相关功能
package admin

import (
	"errors"
	"time"

	"github.com/QOSGroup/qmoon/db"
	"github.com/QOSGroup/qmoon/db/model"
	"github.com/QOSGroup/qmoon/utils"
)

type Account struct {
	ma    *model.Account
	mapps []*model.App

	ID          int64     `json:"id"`
	Mail        string    `json:"mail"`        // mail
	Name        string    `json:"name"`        // name
	Avatar      string    `json:"avatar"`      // avatar
	Description string    `json:"description"` // description
	Status      int64     `json:"status"`      // status
	CreatedAt   time.Time `json:"created_at"`  // created_at
}

type App struct {
	mapp      *model.App
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	SecretKey string `json:"secretKey"`
	Status    int64  `json:"desc"`
	AccountID int64  `json:"accountId"`
}

func covertToAccount(ma *model.Account) *Account {
	return &Account{
		ma:          ma,
		ID:          ma.ID,
		Mail:        ma.Mail.String,
		Name:        ma.Name.String,
		Avatar:      ma.Avatar.String,
		Description: ma.Description.String,
		Status:      ma.Status.Int64,
		CreatedAt:   ma.CreatedAt.Time,
	}
}

func covertToApp(app *model.App) *App {
	return &App{
		mapp:      app,
		ID:        app.ID,
		Name:      app.Name.String,
		SecretKey: app.SecretKey.String,
		Status:    app.Status.Int64,
		AccountID: app.AccountID.Int64,
	}
}

func Accounts(offset, limit int64) ([]*Account, error) {
	mas, err := model.AccountFilter(db.Db, "", offset, limit)
	if err != nil {
		return nil, err
	}

	var as []*Account
	for _, v := range mas {
		as = append(as, covertToAccount(v))
	}

	return as, nil
}

func CreateAccount(mail, password string) (*Account, error) {
	ma := &model.Account{
		Mail:      utils.NullString(mail),
		Password:  utils.NullString(utils.EncryptPwd([]byte(password))),
		CreatedAt: utils.NullTime(time.Now()),
	}

	err := ma.Insert(db.Db)
	if err != nil {
		return nil, err
	}

	return covertToAccount(ma), nil
}

func RetrieveAccountByID(id int64) (*Account, error) {
	ma, err := model.AccountByID(db.Db, id)
	if err != nil {
		return nil, err
	}

	return covertToAccount(ma), nil
}
func RetrieveAccountByMail(mail string) (*Account, error) {
	ma, err := model.AccountByMail(db.Db, utils.NullString(mail))
	if err != nil {
		return nil, err
	}

	return covertToAccount(ma), nil
}

func (a Account) CheckPassword(pwd string) bool {

	return utils.EncryptPwd([]byte(pwd)) == a.ma.Password.String
}

func (a Account) Apps() ([]*App, error) {
	if a.mapps == nil {
		apps, err := model.AppsByAccountID(db.Db, utils.NullInt64(a.ma.ID))
		if err != nil {
			return nil, err
		}

		a.mapps = apps
	}

	var apps []*App
	for _, v := range a.mapps {
		apps = append(apps, covertToApp(v))
	}

	return apps, nil
}

func (a *Account) AppByID(id int64) (*App, error) {
	if a.mapps == nil {
		apps, err := model.AppsByAccountID(db.Db, utils.NullInt64(a.ma.ID))
		if err != nil {
			return nil, err
		}

		a.mapps = apps
	}
	for _, v := range a.mapps {
		if v.ID == id {
			return covertToApp(v), nil
		}
	}

	return nil, errors.New("not found")
}

func (a *Account) DeleteByID(id int64) error {
	if a.mapps == nil {
		apps, err := model.AppsByAccountID(db.Db, utils.NullInt64(a.ma.ID))
		if err != nil {
			return err
		}

		a.mapps = apps
	}
	for k, v := range a.mapps {
		if v.ID == id {
			err := v.Delete(db.Db)
			if err != nil {
				return err
			}
			a.mapps = append(a.mapps[:k], a.mapps[k+1:]...)
			return nil
		}
	}

	return errors.New("not found")
}

func (a *Account) CreateApp(name string) (*App, error) {
	token := utils.MD5([]byte(time.Now().String()))
	mapp := &model.App{
		Name:      utils.NullString(name),
		SecretKey: utils.NullString(token),
		AccountID: utils.NullInt64(a.ma.ID),
		CreatedAt: utils.NullTime(time.Now()),
	}

	err := mapp.Insert(db.Db)
	if err != nil {
		return nil, err
	}

	return covertToApp(mapp), nil
}

func AppBySecretKey(secretKey string) (*App, error) {
	app, err := model.AppBySecretKey(db.Db, utils.NullString(secretKey))
	if err != nil {
		return nil, err
	}

	return covertToApp(app), nil
}

func (app App) Account() (*Account, error) {
	a, err := app.mapp.Account(db.Db)
	if err != nil {
		return nil, err
	}

	return covertToAccount(a), nil
}
