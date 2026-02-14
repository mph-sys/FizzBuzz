package models

type ResponseError struct {
	Errors []string `json:"errors"`
}

type ResponseSuccess struct {
	Data any `json:"data"`
}
