package models

type BrandRequest struct {
	BrandName string `json:"brand_name"`
}

type BrandResponse struct {
	Id        int64  `json:"id"`
	BrandName string `json:"brand_name"`
	Logo      string `json:"logo"`
}

type BrandUpdateReq struct {
	Id        int64  `json:"id"`
	BrandName string `json:"brand_name"`
}

type GetBrandInfo struct {
	Id         int64             `json:"id"`
	BrandName  string            `json:"brand_name"`
	Categories []BrandCategories `json:"categories"`
}

type BrandCategories struct {
	Id           int64  `json:"id"`
	CategoryName string `json:"category_name"`
}

type BrandCategoryReq struct {
	BrandId    int64 `json:"brand_id"`
	CategoryId int64 `json:"category_id"`
}
