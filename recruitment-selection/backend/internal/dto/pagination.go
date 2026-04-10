package dto

// PaginatedResponse wraps any list response with pagination metadata.
type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Pages int64       `json:"pages"`
}

// NewPaginated builds a PaginatedResponse from its parts.
func NewPaginated(data interface{}, total int64, page, limit int) PaginatedResponse {
	pages := total / int64(limit)
	if total%int64(limit) != 0 {
		pages++
	}
	return PaginatedResponse{
		Data:  data,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pages,
	}
}
