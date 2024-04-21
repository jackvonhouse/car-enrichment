package dto

type Car struct {
	ID     int64  `json:"id"`
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Owner  Owner  `json:"owner"`
}

type EnrichmentCar struct {
	Car Car
	Err error
}

type CreateCar struct {
	RegNumbers []string `json:"reg_numbers"`
}

type Pagination struct {
	Offset int
	Limit  int
}

type Filter struct {
	RegNum          string
	Mark            string
	Model           string
	Year            int
	OwnerName       string
	OwnerSurname    string
	OwnerPatronymic string
}
