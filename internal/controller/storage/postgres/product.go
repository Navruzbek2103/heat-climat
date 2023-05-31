package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/db"
)

type ProductRepo struct {
	db *db.Postgres
}

func NewProduct(db *db.Postgres) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (p *ProductRepo) CreateProduct(product *repo.ProductRequest) (*repo.ProductResponse, error) {
	var (
		create, update time.Time
		res            = repo.ProductResponse{}
	)
	query := `
	INSERT INTO products(
		brand_id, category_id, title, 
		description, price, display_type, 
		os_type, camera, dioganal, characteristics) 
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	RETURNING 
		id, brand_id, category_id, title, description, price, 
		display_type, os_type, camera, dioganal, characteristics, rating, created_at, updated_at`

	err := p.db.Pool.QueryRow(context.Background(), query,
		product.BrandId, product.CategoryId,
		product.Title, product.Description,
		product.Price, product.DisplayType,
		product.OsType, product.Camera, product.Dioganal, product.Characteristics).
		Scan(&res.Id, &res.BrandId, &res.CategoryId,
			&res.Title, &res.Description, &res.Price,
			&res.DisplayType, &res.OsType, &res.Camera, &res.Dioganal,
			&res.Characteristics, &res.Rating, &create, &update,
		)
	if err != nil {
		return &repo.ProductResponse{}, err
	}
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)

	media := repo.MediaRes{}
	for _, m := range product.MediaLinks {
		query2 := `
		INSERT INTO
			product_media(product_id, media_link)
		VALUES
			($1, $2)
		RETURNING 
			id, product_id, media_link
		`
		err = p.db.Pool.QueryRow(context.Background(), query2, res.Id, m.MediaLink).
			Scan(&media.Id, &media.ProductId, &media.MediaLink)
		if err != nil {
			return &repo.ProductResponse{}, err
		}
		res.MediaLinks = append(res.MediaLinks, &media)
	}
	return &res, nil
}

func (p *ProductRepo) GetProductById(id int64) (*repo.ProductResponse, error) {
	var (
		create, update time.Time
		res            = repo.ProductResponse{}
	)
	query := `
	SELECT 
		id, category_id, brand_id, title, description, price, 
		display_type, os_type, camera, dioganal, characteristics, rating, created_at, updated_at
	FROM 
		products 
	WHERE 
		id=$1 AND deleted_at IS NULL`
	err := p.db.Pool.QueryRow(context.Background(), query, id).Scan(
		&res.Id, &res.BrandId, &res.CategoryId,
		&res.Title, &res.Description, &res.Price,
		&res.DisplayType, &res.OsType, &res.Camera, &res.Dioganal,
		&res.Characteristics, &res.Rating, &create, &update,
	)
	if err == sql.ErrNoRows {
		return &repo.ProductResponse{}, nil
	}
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)

	rows, err := p.db.Pool.Query(context.Background(),
		`
	SELECT 
		id, product_id, media_link 
	FROM 
		product_media WHERE product_id=$1`, id)
	if err == sql.ErrNoRows {
		return &repo.ProductResponse{}, nil
	}
	for rows.Next() {
		var media repo.MediaRes
		err = rows.Scan(&media.Id, &media.ProductId, &media.MediaLink)
		if err != nil {
			return &repo.ProductResponse{}, err
		}
		res.MediaLinks = append(res.MediaLinks, &media)
	}
	return &res, nil
}

func (p *ProductRepo) UpdateProduct(product *repo.ProductUpdateReq) (*repo.ProductResponse, error) {
	var (
		create_time, update_time time.Time
		res                      = repo.ProductResponse{}
	)
	query := `
	UPDATE 
		products 
	SET 
		title=$1, description=$2, price=$3,
		display_type=$4, os_type=$5, camera=$6, dioganal=$7, 
		characteristics=$8, updated_at=NOW() 
	WHERE 
		id=$9 AND deleted_at IS NULL
	RETURNING 
		id, brand_id, category_id, title, description, 
		price, display_type, os_type, camera, dioganal, 
		characteristics, rating, created_at, updated_at`
	err := p.db.Pool.QueryRow(context.Background(), query,
		product.Title, product.Description, product.Price,
		product.DisplayType, product.OsType, product.Camera,
		product.Dioganal, product.Characteristics, product.Id).
		Scan(&res.Id, &res.BrandId, &res.CategoryId,
			&res.Title, &res.Description, &res.Price,
			&res.DisplayType, &res.OsType, &res.Camera, &res.Dioganal,
			&res.Characteristics, &res.Rating, &create_time, &update_time,
		)
	if err != nil {
		fmt.Println("error while update product info: ", err)
		return &repo.ProductResponse{}, err
	}
	res.CreatedAt = create_time.Format(time.RFC1123)
	res.UpdatedAt = update_time.Format(time.RFC1123)
	var media repo.MediaRes
	for _, m := range product.MediaLink {
		query1 := `
	UPDATE 
		product_media 
	SET 
		media_link=$1 
	WHERE 
		id=$2 
	RETURNING 
		id, product_id, media_link`
		err = p.db.Pool.QueryRow(context.Background(), query1, m.MediaLink, m.Id).
			Scan(&media.Id, &media.ProductId, &media.MediaLink)
		if err != nil {
			fmt.Println("error while update media: ", err)
			return &repo.ProductResponse{}, err
		}
		res.MediaLinks = append(res.MediaLinks, &media)
	}
	return &res, nil
}

func (p *ProductRepo) GetAllProducts(product *repo.AllProductsParams) (*repo.AllProducts, error) {
	var (
		create, update time.Time
		res            = repo.AllProducts{}
	)
	query := `
	SELECT 
		id, brand_id, category_id, title, price, created_at, updated_at
	FROM 
		products 
	WHERE 
		deleted_at IS NULL AND title 
	ILIKE 
		$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
			`
	rows, err := p.db.Pool.Query(context.Background(), query,
		"%"+product.Search+"%", product.Limit, (product.Page-1)*product.Limit)
	if err != nil {
		fmt.Println("Error while get all products")
		return &repo.AllProducts{}, err
	}
	for rows.Next() {
		temp := repo.ProductForList{}
		err = rows.Scan(
			&temp.Id,
			&temp.BrandId,
			&temp.CategoryId,
			&temp.Title,
			&temp.Price, &create, &update)
		if err != nil {
			return &repo.AllProducts{}, err
		}
		query1 := `
		SELECT 
			id, media_link 
		FROM 
			product_media
		WHERE 
			product_id=$1`

		media_rows, err := p.db.Pool.Query(context.Background(), query1, temp.Id)
		if err != nil {
			return &repo.AllProducts{}, err
		}
		for media_rows.Next() {
			image := repo.MediaRes{}
			err = media_rows.Scan(&image.Id, &image.MediaLink)
			if err != nil {
				return &repo.AllProducts{}, err
			}
		}
		res.Products = append(res.Products, &temp)
	}
	return &res, nil
}

func (p *ProductRepo) DeleteProductById(id int64) (*repo.Empty, error) {
	_, err := p.db.Pool.Exec(context.Background(),
		`UPDATE products SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return &repo.Empty{}, err
	}
	_, err = p.db.Pool.Exec(context.Background(), `DELETE FROM product_media WHERE product_id=$1`, id)
	if err != nil {
		return &repo.Empty{}, err
	}
	return &repo.Empty{}, nil
}
