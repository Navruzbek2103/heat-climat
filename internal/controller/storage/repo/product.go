package repo

type ProductRequest struct {
	BrandId         int64
	CategoryId      int64
	Title           string
	Description     string
	Price           float64
	DisplayType     string
	OsType          string
	Camera          string
	Dioganal        float64
	Characteristics string
	MediaLinks      []*MediaReq
}
type MediaReq struct {
	MediaLink string
}

type ProductResponse struct {
	Id              int64
	BrandId         int64
	CategoryId      int64
	Title           string
	Description     string
	Price           float64
	DisplayType     string
	OsType          string
	Camera          string
	Dioganal        float64
	Rating          float64
	Characteristics string
	CreatedAt       string
	UpdatedAt       string
	MediaLinks      []*MediaRes
}
type MediaRes struct {
	Id        int64
	ProductId int64
	MediaLink string
}

type ProductId struct {
	Id int64
}

type AllProductsParams struct {
	Page   int64
	Limit  int64
	Search string
}

type ProductForList struct {
	Id         int64
	BrandId    int64
	CategoryId int64
	Title      string
	Price      float64
	MediaLinks []*MediaRes
}
type AllProducts struct {
	Products []*ProductForList
}

type ProductUpdateReq struct {
	Id              int64
	Title           string
	Description     string
	Price           float64
	DisplayType     string
	OsType          string
	Camera          string
	Dioganal        float64
	Characteristics string
	MediaLink       []*UpdateMediaLink
}

type UpdateMediaLink struct {
	Id        int64
	MediaLink string
}

type ProductStorageI interface {
	CreateProduct(p *ProductRequest) (*ProductResponse, error)
	UpdateProduct(p *ProductUpdateReq) (*ProductResponse, error)
	GetProductById(id int64) (*ProductResponse, error)
	GetAllProducts(params *AllProductsParams) (*AllProducts, error)
	DeleteProductById(id int64) (*Empty, error)
}
