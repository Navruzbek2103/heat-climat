package repo

type AdminRequest struct {
	Id           string
	UserName     string
	Password     string
	RefreshToken string
}

type AdminResponse struct {
	Id           string
	UserName     string
	Password     string
	AccessToken  string
	RefreshToken string
	CreatedAt    string
	UpdatedAt    string
}

type UpdateAdminReq struct {
	Id       string
	UserName string
	Password string
}

type AllAdminParams struct {
	Search string
}

type AllAdmins struct {
	Admins []AdminResponse
}

type CheckFieldReq struct {
	Field string
	Value string
}

type CheckFieldRes struct {
	Exists bool
}

type AdminStorageI interface {
	AddAdmin(a *AdminRequest) (*AdminResponse, error)
	GetAdminInfo(id string) (*AdminResponse, error)
	GetAllAdmins(keyword string) (*AllAdmins, error)
	UpdateAdmin(a *UpdateAdminReq) (*AdminResponse, error)
	DeleteAdmin(id string) (*Empty, error)
	CheckField(a *CheckFieldReq) (*CheckFieldRes, error)
	GetAdmin(username string) (*AdminResponse, error)
}
