package response

type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
