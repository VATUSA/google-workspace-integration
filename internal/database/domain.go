package database

import "gorm.io/gorm"

type Domain struct {
	*gorm.Model
	Facility string `gorm:"size:4"`
	Domain   string `gorm:"size:120"`
}

func domainQuery() *gorm.DB {
	query := DB.Model(&Domain{})
	return query
}

func FetchDomains() ([]Domain, error) {
	var domains []Domain
	result := domainQuery().Find(&domains)
	if result.Error != nil {
		return nil, result.Error
	}
	return domains, nil
}

func FetchDomainsByFacility(facility string) ([]Domain, error) {
	var domains []Domain
	result := domainQuery().Where("facility = ?", facility).Find(&domains)
	if result.Error != nil {
		return nil, result.Error
	}
	return domains, nil
}

func (d *Domain) Save() error {
	result := DB.Save(d)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
