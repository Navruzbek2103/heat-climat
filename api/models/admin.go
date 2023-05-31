package models

type AdminLoginReq struct {
	UserName string `json:"AdminName"`
	Password string `json:"password"`
}
type AdminLoginRes struct {
	Id          string `json:"id"`
	UserName    string `json:"username"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type AdminReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type AdminRes struct {
	Id          string `json:"id"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GetAdminProdile struct {
	Id        string `json:"id"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AllAdmins struct {
	Admins []GetAdminProdile `json:"admins"`
}
