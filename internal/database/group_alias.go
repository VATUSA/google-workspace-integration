package database

type GroupAlias struct {
	Email             string `gorm:"primaryKey;size:120"`
	GroupPrimaryEmail string `gorm:"size:120"`
	Group             *Group
	Facility          string `gorm:"size:12"`
	Domain            string `gorm:"size:120"`
}

func (a *GroupAlias) Save() error {
	result := DB.Save(a)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (a *GroupAlias) Delete() error {
	result := DB.Delete(a)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
