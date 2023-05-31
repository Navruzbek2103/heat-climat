package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/db"
)

type NewsProductRepo struct {
	db *db.Postgres
}

func NewNewsProduct(db *db.Postgres) *NewsProductRepo {
	return &NewsProductRepo{
		db: db,
	}
}

func (n *NewsProductRepo) CreateNewProduct(new *repo.NewProductRequest) (*repo.NewProductResponse, error) {
	var (
		create, update time.Time
		res            = repo.NewProductResponse{}
	)
	query := `
	INSERT INTO new_products(
		brand_id, title, description, new_price,
		old_price, display_type, os_type, camera,
		dioganal, characteristics) 
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING 
		id, 
		brand_id, 
		title, description, new_price,
		old_price, display_type, os_type, camera,
		dioganal, characteristics, created_at, updated_at`
	err := n.db.Pool.QueryRow(context.Background(), query,
		new.BrandId, new.Title, new.Description,
		new.NewPrice, new.OldPrice, new.DisplayType,
		new.OsType, new.Camera, new.Dioganal, new.Charactestics,
	).Scan(
		&res.Id, &res.BrandId, &res.Title, &res.Description, &res.NewPrice,
		&res.OldPrice, &res.DisplayType, &res.OsType,
		&res.Camera, &res.Dioganal, &res.Charactestics, &create, &update)
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	if err != nil {
		return &repo.NewProductResponse{}, err
	}

	var media repo.MediaRes
	for _, m := range new.MediaLinks {
		query1 := `
		INSERT INTO 
			new_product_media(new_product_id, media_link)
		VALUES
			($1, $2)
		RETURNING 
			id, new_product_id, media_link`
		err = n.db.Pool.QueryRow(context.Background(), query1, res.Id, m.MediaLink).Scan(&media.Id, &media.ProductId, &media.MediaLink)
		if err != nil {
			return &repo.NewProductResponse{}, err
		}
		res.MediaLinks = append(res.MediaLinks, &media)
	}

	return &res, nil
}

func (n *NewsProductRepo) GetNewProductById(id int64) (*repo.NewProductResponse, error) {
	var (
		create, update time.Time
		res            = repo.NewProductResponse{}
	)
	query := `
	SELECT 
		id, brand_id, title, description, new_price,
		old_price, display_type, os_type, camera,
		dioganal, characteristics, created_at, updated_at
	FROM 
		new_products 
	WHERE 
		deleted_at IS NULL AND id=$1`
	err := n.db.Pool.QueryRow(context.Background(), query, id).
		Scan(
			&res.Id, &res.BrandId, &res.Title,
			&res.Description, &res.NewPrice, &res.OldPrice,
			&res.DisplayType, &res.OsType, &res.Camera,
			&res.Dioganal, &res.Charactestics, &create, &update,
		)
	if res.Title == "" {
		res.CreatedAt, res.UpdatedAt = "", ""
	} else {
		res.CreatedAt = create.Format(time.RFC1123)
		res.UpdatedAt = update.Format(time.RFC1123)
	}
	if err == sql.ErrNoRows {
		return &repo.NewProductResponse{}, nil
	}
	query2 := `
	SELECT 
		id, new_product_id, media_link
	FROM 
		new_product_media
	WHERE 
		new_product_id=$1`
	rows, err := n.db.Pool.Query(context.Background(), query2, res.Id)
	if err != nil {
		return &repo.NewProductResponse{}, err
	}
	for rows.Next() {
		var media repo.MediaRes
		err = rows.Scan(&media.Id, &media.ProductId, &media.MediaLink)
		if err != nil {
			return &repo.NewProductResponse{}, err
		}
		res.MediaLinks = append(res.MediaLinks, &media)
	}
	return &res, nil
}

func (n *NewsProductRepo) GetAllNewProducts(params *repo.AllNewProductsParams) (*repo.AllNewProducts, error) {
	var (
		res = repo.AllNewProducts{}
	)
	query := `
	SELECT 
		id, brand_id, title, new_price, old_price 
	FROM 
		new_products 
	WHERE 
		deleted_at IS NULL AND title ILIKE $1 ORDER BY created_at DESC
	LIMIT 
		$2 OFFSET $3`
	rows, err := n.db.Pool.Query(context.Background(), query,
		"%"+params.Search+"%",
		params.Limit, params.Page,
	)
	if err != nil {
		return &repo.AllNewProducts{}, err
	}
	for rows.Next() {
		temp := repo.NewProductForList{}
		err = rows.Scan(
			&temp.Id,
			&temp.BrandId,
			&temp.Title,
			&temp.NewPrice,
			&temp.OldPrice,
		)
		if err != nil {
			return &repo.AllNewProducts{}, err
		}

		query2 := `
		SELECT 
			id, media_link
		FROM 
			new_product_media
		WHERE 
			new_product_id=$1
			`
		mediaRows, err := n.db.Pool.Query(context.Background(), query2, temp.Id)
		if err != nil {
			return &repo.AllNewProducts{}, err
		}
		for mediaRows.Next() {
			media_temp := repo.MediaRes{}
			err = mediaRows.Scan(&media_temp.Id, &media_temp.MediaLink)
			if err != nil {
				return &repo.AllNewProducts{}, err
			}
			temp.MediaLinks = append(temp.MediaLinks, &media_temp)
		}
		res.NewProducts = append(res.NewProducts, &temp)
	}
	return &res, nil
}

func (n *NewsProductRepo) UpdateNewProduct(new *repo.NewProductUpdateReq) (*repo.NewProductResponse, error) {
	var (
		create, update time.Time
		res            = repo.NewProductResponse{}
	)
	query := `
	UPDATE 
		new_products 
	SET 
		title=$1, description=$2, new_price=$3,
		old_price=$4, display_type=$5, os_type=$6, camera=$7,
		dioganal=$8, characteristics=$9,
		updated_at=NOW()
	WHERE 
		id=$10 AND deleted_at IS NULL
	RETURNING 
		id, brand_id, title, description, new_price,
		old_price, display_type, os_type, camera,
		dioganal, characteristics, created_at, updated_at`
	err := n.db.Pool.QueryRow(context.Background(), query,
		new.Title, new.Description, new.NewPrice,
		new.OldPrice, new.DisplayType, new.OsType,
		new.Camera, new.Dioganal, new.Charactestics, new.Id,
	).Scan(
		&res.Id, &res.BrandId, &res.Title,
		&res.Description, &res.NewPrice,
		&res.OldPrice, &res.DisplayType, &res.OsType,
		&res.Camera, &res.Dioganal, &res.Charactestics, &create, &update)
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	if err != nil {
		return &repo.NewProductResponse{}, err
	}
	for _, m := range new.MediaLink {
		media := repo.MediaRes{}
		query2 := `
		UPDATE 
			new_product_media
		SET 
			media_link=$1
		WHERE 
			id=$2
		RETURNING
			id, new_product_id, media_link`
		err = n.db.Pool.QueryRow(context.Background(), query2, m.MediaLink, m.Id).
			Scan(&media.Id, &media.ProductId, &media.MediaLink)
		if err != nil {
			return &repo.NewProductResponse{}, err
		}
		res.MediaLinks = append(res.MediaLinks, &media)
	}
	return &res, nil
}

func (n *NewsProductRepo) DeleteNewProductById(id int64) (*repo.Empty, error) {
	_, err := n.db.Pool.Exec(context.Background(),
		`UPDATE new_products SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		fmt.Println("error while delete news ", err)
		return &repo.Empty{}, err
	}
	return &repo.Empty{}, nil
}
