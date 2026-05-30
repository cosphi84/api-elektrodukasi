package dto

type PaginationQuery struct {
	Page    int `form:"page"    json:"page"`
	PerPage int `form:"per_page" json:"per_page"`
}

func (p *PaginationQuery) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 10
	}
	if p.PerPage > 100 {
		p.PerPage = 100 // hard cap — no one needs 10000 rows 😅
	}
}

func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.PerPage
}

type PaginatedResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginatedResult[T any](data []T, total int64, query PaginationQuery) PaginatedResult[T] {
	totalPages := int(total) / query.PerPage
	if int(total)%query.PerPage > 0 {
		totalPages++
	}
	return PaginatedResult[T]{
		Data:       data,
		Total:      total,
		Page:       query.Page,
		PerPage:    query.PerPage,
		TotalPages: totalPages,
	}
}
