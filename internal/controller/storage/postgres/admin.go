package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/db"
)

type AdminRepo struct {
	db *db.Postgres
}

func NewAdmin(db *db.Postgres) *AdminRepo {
	return &AdminRepo{
		db: db,
	}
}

func (a *AdminRepo) AddAdmin(admin *repo.AdminRequest) (*repo.AdminResponse, error) {
	var (
		res            repo.AdminResponse
		create, update time.Time
	)

	query := `
	INSERT INTO 
		admins(id, user_name, password, refresh_token)
	VALUES
		($1, $2, $3, $4)
	RETURNING
		id, user_name, password, refresh_token, created_at, updated_at`
	err := a.db.Pool.QueryRow(context.Background(),
		query, admin.Id, admin.UserName, admin.Password, admin.RefreshToken).
		Scan(&res.Id, &res.UserName, &res.Password, &res.RefreshToken, &create, &update)
	if err != nil {
		return &repo.AdminResponse{}, err
	}
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	return &res, nil
}

func (a *AdminRepo) GetAdmin(username string) (*repo.AdminResponse, error) {
	var (
		create, update time.Time
		res            = repo.AdminResponse{}
	)
	query := `
	SELECT 
		id, user_name, password, created_at, updated_at
	FROM 
		admins
	WHERE
		user_name=$1 AND deleted_at IS NULL	
		`
	err := a.db.Pool.QueryRow(context.Background(), query, username).
		Scan(&res.Id, &res.UserName, &res.Password, &create, &update)
	if err != nil {
		fmt.Println("error while getting admin info>> ", err)
		return &repo.AdminResponse{}, err
	}
	if res.UserName == "" {
		res.CreatedAt = ""
		res.UpdatedAt = ""
	} else {
		res.CreatedAt = create.Format(time.RFC1123)
		res.UpdatedAt = update.Format(time.RFC1123)
	}
	return &res, nil
}

func (a *AdminRepo) GetAllAdmins(keyword string) (*repo.AllAdmins, error) {
	res := repo.AllAdmins{}
	var create, update time.Time
	query := `
	SELECT 
		id, user_name, password, created_at, updated_at
	FROM 
		admins
	WHERE 
		user_name ILIKE $1 AND deleted_at IS NULL`
	rows, err := a.db.Pool.Query(context.Background(), query, "%"+keyword+"%")
	if err != nil {
		return &repo.AllAdmins{}, err
	}
	for rows.Next() {
		temp := repo.AdminResponse{}
		err = rows.Scan(&temp.Id, &temp.UserName, &temp.Password, &create, &update)
		if err != nil {
			return &repo.AllAdmins{}, err
		}
		temp.CreatedAt = create.Format(time.RFC1123)
		temp.UpdatedAt = update.Format(time.RFC1123)
		res.Admins = append(res.Admins, temp)
	}
	return &res, nil
}

func (a *AdminRepo) GetAdminInfo(id string) (*repo.AdminResponse, error) {
	var (
		create, update time.Time
		res            repo.AdminResponse
	)
	query := `
	SELECT 
		id, user_name, password, created_at, updated_at
	FROM 
		admins
	WHERE 
		id=$1 AND deleted_at IS NULL`

	err := a.db.Pool.QueryRow(context.Background(), query, id).
		Scan(&res.Id, &res.UserName, &res.Password, &create, &update)
	if err != nil {
		return &repo.AdminResponse{}, err
	}
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	return &res, nil
}

func (a *AdminRepo) UpdateAdmin(admin *repo.UpdateAdminReq) (*repo.AdminResponse, error) {
	var (
		res            repo.AdminResponse
		create, update time.Time
	)
	query := `
	UPDATE 
		admins
	SET
		updated_at=NOW(), user_name=$1, password=$2
	WHERE 
		id=$3 AND deleted_at IS NULL
	RETURNING 
		id, user_name, password, created_at, updated_at`
	err := a.db.Pool.QueryRow(context.Background(),
		query, admin.UserName, admin.Password, admin.Id).
		Scan(&res.Id, &res.UserName, &res.Password, &create, &update)
	if err != nil {
		return &repo.AdminResponse{}, err
	}
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	return &res, nil
}

func (a *AdminRepo) DeleteAdmin(id string) (*repo.Empty, error) {
	_, err := a.db.Pool.Exec(context.Background(), `DELETE FROM admins WHERE id=$1`, id)
	if err != nil {
		return &repo.Empty{}, err
	}
	return &repo.Empty{}, nil
}

func (a *AdminRepo) CheckField(req *repo.CheckFieldReq) (*repo.CheckFieldRes, error) {
	res := &repo.CheckFieldRes{Exists: false}
	query := fmt.Sprintf("SELECT 1 FROM admins WHERE %s=$1", req.Field)
	var temp = 0
	err := a.db.Pool.QueryRow(context.Background(), query, req.Value).Scan(&temp)
	if err == sql.ErrNoRows {
		return res, nil
	} else if err != nil {
		return res, nil
	}
	if temp == 1 {
		res.Exists = true
		return res, nil
	}
	return res, nil
}
