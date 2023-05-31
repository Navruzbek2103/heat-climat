package models

type SuccessInfo struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type FailureInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   error  `json:"error"`
}
