package models

type CompanyReq struct {
	Name         string         `json:"name"`
	Logo         string         `json:"logo"`
	Email        string         `json:"email"`
	PhoneNumbers []string       `json:"phone_numbers"`
	Bio          string         `json:"bio"`
	Addresses    AddressReq     `json:"addresses"`
	SocialMedia  SocialMediaReq `json:"social_medias"`
}

type AddressReq struct {
	District string `json:"district"`
	Street   string `json:"street"`
}

type SocialMediaReq struct {
	Facebook  string `json:"facebook"`
	Instagram string `json:"instagram"`
	Telegram  string `json:"telegram"`
	Twitter   string `json:"twitter"`
}
