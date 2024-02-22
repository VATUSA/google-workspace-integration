package database

type AliasType int

const (
	AliasType_FacilityName AliasType = iota
	AliasType_FacilityPosition
)

type Alias struct {
	Email     string `gorm:"primaryKey;size:120"`
	AccountId uint64
	Account   *Account
	AliasType AliasType
	Facility  string `gorm:"size:12"`
	Role      string `gorm:"size:12"`
}

func (a *Alias) Save() error {
	result := DB.Save(a)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *Alias) Delete() error {
	result := DB.Delete(a)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
