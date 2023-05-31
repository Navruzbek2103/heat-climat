package models

type CategoryReq struct {
	CategoryName string `json:"category_name"`
}

type CategoryUpdateReq struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
}
