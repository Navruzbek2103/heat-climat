package storage

import (
	"gitlab.com/climate.uz/internal/controller/storage/postgres"
	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/db"
)

type StorageI interface {
	Admin() repo.AdminStorageI
	Brand() repo.BrandStorageI
	Category() repo.CategoryStorageI
	Product() repo.ProductStorageI
	NewProduct() repo.NewsStorageI
	User() repo.UserStorageI
	BrandCategory() repo.BrandCategoryStorageI
}

type StoragePg struct {
	Adminrepo         repo.AdminStorageI
	Brandrepo         repo.BrandStorageI
	Categoryrepo      repo.CategoryStorageI
	Productrepo       repo.ProductStorageI
	Userrepo          repo.UserStorageI
	NewsProductrepo   repo.NewsStorageI
	BrandCategoryrepo repo.BrandCategoryStorageI
}

func NewStoragePg(db *db.Postgres) StorageI {
	return &StoragePg{
		Adminrepo:         postgres.NewAdmin(db),
		Brandrepo:         postgres.NewBrand(db),
		Categoryrepo:      postgres.NewCategory(db),
		Productrepo:       postgres.NewProduct(db),
		NewsProductrepo:   postgres.NewNewsProduct(db),
		Userrepo:          postgres.NewUser(db),
		BrandCategoryrepo: postgres.NewBrandCategory(db),
	}
}

func (a *StoragePg) Admin() repo.AdminStorageI {
	return a.Adminrepo
}

func (b *StoragePg) Brand() repo.BrandStorageI {
	return b.Brandrepo
}
func (s *StoragePg) Category() repo.CategoryStorageI {
	return s.Categoryrepo
}

func (s *StoragePg) Product() repo.ProductStorageI {
	return s.Productrepo
}

func (s *StoragePg) NewProduct() repo.NewsStorageI {
	return s.NewsProductrepo
}

func (s *StoragePg) User() repo.UserStorageI {
	return s.Userrepo
}

func (bc *StoragePg) BrandCategory() repo.BrandCategoryStorageI {
	return bc.BrandCategoryrepo
}
