package domain

type Pages struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type Range struct {
	Start int
	End   int
}
