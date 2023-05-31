package postgres

import (
	"context"
	"database/sql"
	"time"

	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/db"
)

type UserRepo struct {
	db *db.Postgres
}

func NewUser(db *db.Postgres) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) CreateUser(user *repo.UserRequest) (*repo.UserResponse, error) {
	var (
		create, update time.Time
		res            = repo.UserResponse{}
	)
	query := `INSERT INTO 
		users(phone_number) VALUES($1) 
		RETURNING id, phone_number, created_at, updated_at`

	err := u.db.Pool.QueryRow(context.Background(),
		query,
		user.PhoneNumber).
		Scan(&res.Id, &res.PhoneNumber, &create, &update)
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	if err != nil {
		return &repo.UserResponse{}, err
	}

	return &res, nil
}

func (u *UserRepo) GetUserById(id int64) (*repo.UserResponse, error) {
	var (
		create, update time.Time
		res            = repo.UserResponse{}
	)
	query := `SELECT 
			id, phone_number, 
			created_at, updated_at 
		FROM users WHERE deleted_at IS NULL AND id=$1`
	err := u.db.Pool.QueryRow(context.Background(), query, id).Scan(
		&res.Id, &res.PhoneNumber, &create, &update,
	)
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	if err == sql.ErrNoRows {
		return &repo.UserResponse{}, nil
	}

	return &res, nil
}

func (u *UserRepo) GetAllUser(params *repo.AllUsersParams) (*repo.AllUsers, error) {
	var (
		create, update time.Time
		res            repo.AllUsers
	)
	query := `SELECT 
			id, phone_number, created_at, updated_at 
			FROM users WHERE deleted_at IS NULL AND phone_number ILIKE $1 LIMIT $2 OFFSET $3`
	rows, err := u.db.Pool.Query(context.Background(), query, "%"+params.Search+"%", params.Limit, (params.Page-1)*params.Limit)
	if err != nil {
		return &repo.AllUsers{}, err
	}
	for rows.Next() {
		response := repo.UserResponse{}
		err = rows.Scan(
			&response.Id,
			&response.PhoneNumber,
			&create,
			&update,
		)
		response.CreatedAt = create.Format(time.RFC1123)
		response.UpdatedAt = update.Format(time.RFC1123)
		if err != nil {
			return &repo.AllUsers{}, err
		}
		res.Users = append(res.Users, &response)
	}

	return &res, nil
}

func (u *UserRepo) UpdateUser(user *repo.UserUpdateReq) (*repo.UserResponse, error) {
	var (
		create, update time.Time
		res            = repo.UserResponse{}
	)
	query := `UPDATE users 
		SET phone_number=$1, updated_at=NOW() 
		WHERE id=$2 AND deleted_at IS NULL
		RETURNING id, phone_number, created_at, updated_at`
	err := u.db.Pool.QueryRow(context.Background(), query,
		user.PhoneNumber, user.Id).Scan(
		&res.Id, &res.PhoneNumber, &create, &update,
	)
	res.CreatedAt = create.Format(time.RFC1123)
	res.UpdatedAt = update.Format(time.RFC1123)
	if err != nil {
		return &repo.UserResponse{}, err
	}

	return &res, nil
}

func (u *UserRepo) DeleteUser(id int64) (*repo.Empty, error) {
	_, err := u.db.Pool.Exec(context.Background(), `UPDATE users SET deleted_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return &repo.Empty{}, err
	}
	return &repo.Empty{}, nil
}
