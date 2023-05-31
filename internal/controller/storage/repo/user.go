package repo

type UserRequest struct {
	PhoneNumber string
}

type UserResponse struct {
	Id          int64
	PhoneNumber string
	CreatedAt   string
	UpdatedAt   string
}

type UserId struct {
	Id int64
}

type UserUpdateReq struct {
	Id          int64
	PhoneNumber string
}

type AllUsersParams struct {
	Page   int64
	Limit  int64
	Search string
}

type AllUsers struct {
	Users []*UserResponse
}

type UserStorageI interface {
	CreateUser(u *UserRequest) (*UserResponse, error)
	UpdateUser(u *UserUpdateReq) (*UserResponse, error)
	GetUserById(id int64) (*UserResponse, error)
	GetAllUser(params *AllUsersParams) (*AllUsers, error)
	DeleteUser(id int64) (*Empty, error)
}
