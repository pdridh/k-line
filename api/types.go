package api

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func CalculatePaginationMeta(total, page, limit int) *PaginationMeta {
	if limit == 0 {
		return nil
	}

	totalPages := (total + limit - 1) / limit

	return &PaginationMeta{
		Total:      total,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}
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

type SuccessResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewSuccessResponse(status int, code string, message string, data any) *SuccessResponse {
	return &SuccessResponse{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
	}
}
