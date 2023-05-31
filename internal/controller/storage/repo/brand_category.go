package repo

type BrandCategoryReq struct {
	BrandId    int64
	CategoryId int64
}

type BrandCategoryStorageI interface {
	CreateBrandCategory(*BrandCategoryReq) error
}



