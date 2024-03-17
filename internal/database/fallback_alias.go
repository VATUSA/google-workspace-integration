package database

type FallbackAlias struct {
	Email    string `gorm:"primaryKey;size:120"`
	Facility string `gorm:"size:12"`
	Role     string `gorm:"size:12"`
}

func FetchFallbackAliases() ([]FallbackAlias, error) {
	var aliases []FallbackAlias
	result := DB.Model(&FallbackAlias{}).Find(&aliases)
	if result.Error != nil {
		return nil, result.Error
	}
	return aliases, nil
}

func (a *FallbackAlias) Save() error {
	result := DB.Save(a)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *FallbackAlias) Delete() error {
	result := DB.Delete(a)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
