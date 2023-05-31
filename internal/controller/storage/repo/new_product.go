package repo

type NewProductRequest struct {
	BrandId       int64
	Title         string
	Description   string
	NewPrice      float64
	OldPrice      float64
	DisplayType   string
	OsType        string
	Camera        string
	Dioganal      float64
	Charactestics string
	MediaLinks    []*MediaReq
}

type NewProductResponse struct {
	Id            int64
	BrandId       int64
	Title         string
	Description   string
	NewPrice      float64
	OldPrice      float64
	DisplayType   string
	OsType        string
	Camera        string
	Dioganal      float64
	Charactestics string
	CreatedAt     string
	UpdatedAt     string
	MediaLinks    []*MediaRes
}

type NewProductId struct {
	Id int64
}

type NewProductUpdateReq struct {
	Id            int64
	BrandId       int64
	Title         string
	Description   string
	NewPrice      float64
	OldPrice      float64
	DisplayType   string
	OsType        string
	Camera        string
	Dioganal      float64
	Charactestics string
	MediaLink     []*UpdateMediaLink
}

type AllNewProductsParams struct {
	Page   int64
	Limit  int64
	Search string
}

type NewProductForList struct {
	Id         int64
	BrandId    int64
	Title      string
	NewPrice   float64
	OldPrice   float64
	MediaLinks []*MediaRes
}
type AllNewProducts struct {
	NewProducts []*NewProductForList
}

type NewsStorageI interface {
	CreateNewProduct(n *NewProductRequest) (*NewProductResponse, error)
	UpdateNewProduct(n *NewProductUpdateReq) (*NewProductResponse, error)
	GetNewProductById(id int64) (*NewProductResponse, error)
	GetAllNewProducts(params *AllNewProductsParams) (*AllNewProducts, error)
	DeleteNewProductById(id int64) (*Empty, error)
}
