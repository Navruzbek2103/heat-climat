package models

type UserReq struct {
	PhoneNumber string `json:"phone_number" example:"998"`
}

type UserRes struct {
	Id          int64  `json:"id"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   string `json:"created_at"`
	UpdateAt    string `json:"updated_at"`
}

type UpdateUserReq struct {
	Id          int64  `json:"id"`
	PhoneNumber string `json:"phone_number"`
}
