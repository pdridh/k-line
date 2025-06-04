package menu

type MenuFilters struct {
	Search string `json:"search"`
	Page   int32  `json:"page"` // The request is sent as page but converted to offset for db
	Limit  int32  `json:"limit"`
}

func (f *MenuFilters) Validate(maxLimit int32, defaultLimit int32) {
	// Normalize pagination
	if f.Limit <= 0 || f.Limit > maxLimit {
		f.Limit = defaultLimit
	}

	if f.Page < 1 {
		f.Page = 1
	}

	if f.Search != "" {
		f.Search = "%" + f.Search + "%"
	}
}
