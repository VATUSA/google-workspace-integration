package database

type Group struct {
	PrimaryEmail string `gorm:"primaryKey;size:120"`
	DisplayName  string
	Facility     string
	GroupType    string
}

func FetchGroups() ([]Group, error) {
	var groups []Group
	result := DB.Model(&Group{}).Find(&groups)
	if result.Error != nil {
		return nil, result.Error
	}
	return groups, nil
}

func (g *Group) Save() error {
	result := DB.Save(g)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (g *Group) Delete() error {
	result := DB.Delete(g)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
