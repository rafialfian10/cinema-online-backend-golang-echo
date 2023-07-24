package dto

type SuccessResult struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorResult struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
