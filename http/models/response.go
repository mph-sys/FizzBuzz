package models

type ResponseError struct {
	Errors []string `json:"errors"`
}
