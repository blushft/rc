package rc

type Pagination struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type Status struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`

	Status  string `json:"status"`
	Message string `json:"message"`
}
