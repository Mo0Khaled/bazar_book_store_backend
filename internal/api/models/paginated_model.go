package models

import "math"

type Paginated struct {
	CurrentPage int `json:"current_page"`
	NextPage    int `json:"next_page,omitempty"`
	PrevPage    int `json:"prev_page,omitempty"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
}

func ToPaginated(page int, limit int, totalCount int) Paginated {
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	response := Paginated{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  totalCount,
	}
	if page < totalPages {
		response.NextPage = page + 1
	}

	if page > 1 {
		response.PrevPage = page - 1
	}
	return response
}
