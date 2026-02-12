package request

import "strconv"

type PageLimit struct {
	Page  int
	Limit int
}

func ParsePageLimit(pageStr, limitStr string) PageLimit {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 200 {
		limit = 20
	}
	return PageLimit{Page: page, Limit: limit}
}
