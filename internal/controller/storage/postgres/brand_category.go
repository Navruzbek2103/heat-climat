package postgres

import (
	"context"

	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/db"
)

type BrandCategoryRepo struct {
	db *db.Postgres
}

func NewBrandCategory(db *db.Postgres) *BrandCategoryRepo {
	return &BrandCategoryRepo{
		db: db,
	}
}

func (bc *BrandCategoryRepo) CreateBrandCategory(brandc *repo.BrandCategoryReq) error {
	_, err := bc.db.Pool.Exec(context.Background(), `
	INSERT INTO 
		brand_category(brand_id, category_id)
	VALUES
		($1, $2)`, brandc.BrandId, brandc.CategoryId)
	if err != nil {
		return err
	}
	return nil
}
