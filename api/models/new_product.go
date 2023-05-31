package models

type NewProductReq struct {
	BrandId       int64             `json:"brand_id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	NewPrice      float64           `json:"new_price"`
	OldPrice      float64           `json:"old_price"`
	DisplayType   string            `json:"display_type"`
	OsType        string            `json:"os_type"`
	Camera        string            `json:"camera"`
	Dioganal      float64           `json:"dioganal"`
	Charactestics map[string]string `json:"characteristics" example:"key:value"`
	MediaLinks    []MediaLink       `json:"media links"`
}

type NewProductRes struct {
	Id            int32             `json:"id"`
	BrandId       int64             `json:"brand_id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	NewPrice      float64           `json:"new_price"`
	OldPrice      float64           `json:"old_price"`
	DisplayType   string            `json:"display_type"`
	OsType        string            `json:"os_type"`
	Camera        string            `json:"camera"`
	Dioganal      float64           `json:"dioganal"`
	Charactestics map[string]string `json:"characteristics" example:"key:value"`
	MediaLinks    []MediaLink       `json:"media links"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

type NewProductsList struct {
	Id       int64   `json:"id"`
	BrandId  int64   `json:"brand_id"`
	Title    string  `json:"title"`
	NewPrice float64 `json:"new_price"`
	OldPrice float64 `json:"old_price"`
}
type AllNewProducts struct {
	NewsProducts []NewProductsList `json:"news_products"`
}

type UpdateNewProductReq struct {
	Id            int32             `json:"id"`
	BrandId       int64             `json:"brand_id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	NewPrice      float64           `json:"new_price"`
	OldPrice      float64           `json:"old_price"`
	DisplayType   string            `json:"display_type"`
	OsType        string            `json:"os_type"`
	Camera        string            `json:"camera"`
	Dioganal      float64           `json:"dioganal"`
	Charactestics map[string]string `json:"characteristics" example:"key:value"`
	MediaLinks    []UpdateMediaLink `json:"media_link"`
}
