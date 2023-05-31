package repo

type CategoryRequest struct {
	CategoryName string
}

type CategoryResponse struct {
	Id            int
	CatergoryName string
	CreatedAt     string
	UpdatedAt     string
	Products      []*ProductForList
}

type CategoryId struct {
	Id int
}
type CategoryUpdateReq struct {
	Id           int
	CategoryName string
}
type AllCategoriesParams struct {
	Page   int64
	Limit  int64
	Search string
}
type CategoryList struct {
	Id           int64
	CategoryName string
}
type AllCategory struct {
	Categories []*CategoryList
}
type Empty struct{}

type CategoryRes struct {
	Id           int64
	CategoryName string
}

type CategoryStorageI interface {
	CreateCategory(c *CategoryRequest) (*CategoryResponse, error)
	GetCategoryById(id int) (*CategoryResponse, error)
	GetAllCategories(params *AllCategoriesParams) (*AllCategory, error)
	UpdateCategory(c *CategoryUpdateReq) (*CategoryResponse, error)
	DeleteCategoryById(id int) (*Empty, error)
}
