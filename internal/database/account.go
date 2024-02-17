package database

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	Id               uint64 `gorm:"primaryKey"`
	CID              uint64 `gorm:"column:cid"`
	FirstName        string `gorm:"size:120"`
	LastName         string `gorm:"size:120"`
	PrimaryAlias     string `gorm:"size:120"`
	PrimaryEmail     string `gorm:"size:120"`
	IsManaged        bool
	IsSuspended      bool
	SuspendedAt      *time.Time
	GroupMemberships []GroupMembership
	Aliases          []Alias
}

func accountQuery() *gorm.DB {
	query := DB.Model(&Account{}).
		Preload("GroupMemberships").
		Preload("GroupMemberships.Group").
		Preload("Aliases")
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

func (a *Account) Delete() error {
	result := DB.Delete(a)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
