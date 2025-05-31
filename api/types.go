package api

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse[T any] struct {
	Data []T             `json:"data"`
	Meta *PaginationMeta `json:"meta,omitempty"`
}

func NewPaginatedResponse[T any](data []T, meta *PaginationMeta) PaginatedResponse[T] {
	res := PaginatedResponse[T]{}
	if data == nil {
		data = []T{}
	}

	res.Data = data
	res.Meta = meta

	return res
}
