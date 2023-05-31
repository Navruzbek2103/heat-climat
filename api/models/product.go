package models

type ProductReq struct {
	BrandId         int64             `json:"brand_id"`
	CategoryId      int64             `json:"category_id"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Price           float64           `json:"price"`
	DisplayType     string            `json:"display_type"`
	OsType          string            `json:"os_type"`
	Camera          string            `json:"camera"`
	Dioganal        float64           `json:"dioganal"`
	Characteristics map[string]string `json:"characteristics" example:"key:value"`
	MediaLinks      []MediaLink       `json:"media_link"`
}

type MediaLink struct {
	MediaLink string
}

type UpdateProductReq struct {
	Id              int64             `json:"id"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Price           float64           `json:"price"`
	DisplayType     string            `json:"display_type"`
	OsType          string            `json:"os_type"`
	Camera          string            `json:"camera"`
	Dioganal        float64           `json:"dioganal"`
	Characteristics map[string]string `json:"characteristics" example:"key:value"`
	MediaLinks      []UpdateMediaLink `json:"media_link"`
}
type UpdateMediaLink struct {
	Id        int64  `json:"id"`
	MediaLink string `json:"media_link"`
}

type ProductRes struct {
	Id              int64             `json:"id"`
	BrandId         int64             `json:"brand_id"`
	CategoryId      int64             `json:"category_id"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Price           float64           `json:"price"`
	DisplayType     string            `json:"display_type"`
	OsType          string            `json:"os_type"`
	Camera          string            `json:"camera"`
	Dioganal        float64           `json:"dioganal"`
	Characteristics map[string]string `json:"characteristics" example:"key:value"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
	MediaLinks      []MediaLink       `json:"media_link"`
}

type AllProducts struct {
	Products []ProductRes `json:"products"`
}
