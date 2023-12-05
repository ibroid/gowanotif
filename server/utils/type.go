package utils

type MessagePostBody struct {
	Number     string `json:"number" form:"number" validate:"required,min=5,max=20"`
	Message    string `json:"message" form:"message" validate:"required"`
	ClientName string `json:"client_name" form:"client_name" validate:"required,min=5,max=20"`
}

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

type WaClient struct {
	ClientName        string `json:"client_name" form:"client_name" validate:"required,min=5,max=20"`
	AutoReplyServices bool   `json:"auto_reply" form:"auto_reply" validate:"required"`
}
