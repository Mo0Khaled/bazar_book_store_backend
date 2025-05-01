package helpers

import (
	"net/http"
	"strconv"
)

// GetPaginationFromRequest / It will return page, limit, offset
func GetPaginationFromRequest(r *http.Request) (int, int, int) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	return page, limit, offset
}
