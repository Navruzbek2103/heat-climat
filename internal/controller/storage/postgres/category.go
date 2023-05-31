package postgres

import (
	"context"
	"database/sql"
	"time"

	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/db"
)

type CategoryRepo struct {
	db *db.Postgres
}

func NewCategory(db *db.Postgres) *CategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

func (c *CategoryRepo) CreateCategory(category *repo.CategoryRequest) (*repo.CategoryResponse, error) {
	var (
		res            = repo.CategoryResponse{}
		create, update time.Time
	)
	query, _, err := c.db.Builder.Insert("categories").
		Columns("category_name").Values(category.CategoryName).
		Suffix("RETURNING id, category_name, created_at, updated_at").ToSql()
	if err != nil {
		return &repo.CategoryResponse{}, err
	}
	err = c.db.Pool.QueryRow(context.Background(),
		query, category.CategoryName).Scan(
		&res.Id, &res.CatergoryName,
		&create, &update,
	)
	if err != nil {
		return &repo.CategoryResponse{}, err
	}
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	return &res, nil
}

func (c *CategoryRepo) GetCategoryById(id int) (*repo.CategoryResponse, error) {
	res := repo.CategoryResponse{}
	var create, update time.Time
	query, _, err := c.db.Builder.
		Select(
			"id",
			"category_name",
			"created_at",
			"updated_at").
		From("categories").Where("id=$1 AND deleted_at IS NULL", id).ToSql()
	if err != nil {
		return &repo.CategoryResponse{}, err
	}
	err = c.db.Pool.QueryRow(context.Background(), query, id).
		Scan(&res.Id, &res.CatergoryName, &create, &update)
	if err == sql.ErrNoRows {
		return &repo.CategoryResponse{}, nil
	}
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	query1 := `
	SELECT 
		id, category_id, title, price
	FROM 
		products 
	WHERE
		deleted_at IS NULL AND  cagegory_id=$1`
	rows, err := c.db.Pool.Query(context.Background(), query1, res.Id)
	if err != nil {
		return &repo.CategoryResponse{}, err
	}
	var temp repo.ProductForList
	for rows.Next() {
		err = rows.Scan(&temp.Id, &temp.CategoryId, &temp.Title, &temp.Price)
		if err != nil {
			return &repo.CategoryResponse{}, err
		}

		medias, err := c.db.Pool.Query(context.Background(), 
		`SELECT id, product_id, media_link FROM product_id=$1`, temp.Id)
		if err != nil {
			return &repo.CategoryResponse{}, err
		}
		for medias.Next() {
			media := repo.MediaRes{}
			err = medias.Scan(&media.Id, &media.ProductId, &media.MediaLink)
			if err != nil {
				return &repo.CategoryResponse{}, err
			}
			temp.MediaLinks = append(temp.MediaLinks, &media)
		}
		res.Products = append(res.Products, &temp)
	}

	return &res, nil
}

func (c *CategoryRepo) UpdateCategory(category *repo.CategoryUpdateReq) (*repo.CategoryResponse, error) {
	res := repo.CategoryResponse{}
	var create, update time.Time
	query := `
	UPDATE 
		categories SET category_name=$1, 
		updated_at=NOW() WHERE id=$2
	RETURNING 
		id, category_name, created_at, updated_at`

	err := c.db.Pool.QueryRow(context.Background(), query, category.CategoryName, category.Id).
		Scan(&res.Id, &res.CatergoryName, &create, &update)
	if err != nil {
		return &repo.CategoryResponse{}, err
	}
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	return &res, nil
}

func (c *CategoryRepo) GetAllCategories(param *repo.AllCategoriesParams) (*repo.AllCategory, error) {
	res := repo.AllCategory{}
	query := `
	SELECT
		id, category_name
	FROM 
		categories
	WHERE 
		category_name ILIKE $1 AND deleted_at IS NULL
	LIMIT 
		$2 OFFSET $3`
	rows, err := c.db.Pool.Query(context.Background(),
		query, "%"+param.Search+"%", param.Limit, (param.Page-1)*param.Limit)
	if err != nil {
		return &repo.AllCategory{}, err
	}
	for rows.Next() {
		temp := repo.CategoryList{}
		err = rows.Scan(&temp.Id, &temp.CategoryName)
		if err != nil {
			return &repo.AllCategory{}, err
		}
		res.Categories = append(res.Categories, &temp)
	}
	return &res, nil
}

func (c *CategoryRepo) DeleteCategoryById(id int) (*repo.Empty, error) {
	res := repo.Empty{}
	_, err := c.db.Pool.Exec(context.Background(),
		`UPDATE categories SET deleted_at=NOW() 
	WHERE id=$1 and deleted_at IS NULL`, id)
	if err != nil {
		return &repo.Empty{}, err
	}
	_, err = c.db.Pool.Exec(context.Background(),
		`
	UPDATE products SET deleted_at=NOW()
	WHERE category_id=$1 AND deleted_at IS NULL
	`, id)
	if err != nil {
		return &repo.Empty{}, err
	}
	return &res, nil
}
