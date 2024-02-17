package database

type GroupMembership struct {
	Id         uint64
	GroupEmail string `gorm:"size:120"`
	Group      *Group `gorm:"foreignKey:GroupEmail"`
	AccountID  uint64
	Account    *Account
	IsManager  bool
}

func (m *GroupMembership) Save() error {
	result := DB.Save(m)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *GroupMembership) Delete() error {
	result := DB.Delete(m)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
