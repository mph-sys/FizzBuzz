package models

type FizzBuzzParams struct {
	Int1, Int2, Limit int
	Str1, Str2        string
}

type FizzBuzzStats struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
	Hits  int    `json:"hits"`
}
