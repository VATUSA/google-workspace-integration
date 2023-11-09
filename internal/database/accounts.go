package database

import (
	"fmt"
	"gorm.io/gorm"
)

type Account struct {
	*gorm.Model
	CID          uint `gorm:"column:cid"`
	Alias        string
	PrimaryEmail string
	IsActive     bool
	IsCreated    bool
	IsDeleted    bool
}

func NewAccount(cid uint, alias string) Account {
	return Account{
		CID:          cid,
		Alias:        alias,
		PrimaryEmail: fmt.Sprintf("%s@vatusa.net", alias),
		IsActive:     true,
		IsCreated:    false,
		IsDeleted:    false,
	}
}

func accountQuery() *gorm.DB {
	query := DB.Model(&Account{})
	return query
}

func FetchAccounts() ([]Account, error) {
	var accounts []Account
	result := accountQuery().Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return accounts, nil
}

func (a *Account) Save() error {
	result := DB.Save(a)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
