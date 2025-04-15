package model

type WebResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Pagination struct {
	Limit      int   `json:"limit"`
	Page       int   `json:"page"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type DataWithPagination struct {
	Items  interface{} `json:"items"`
	Paging Pagination  `json:"paging"`
}

type PaginationFilter struct {
	Limit  int    `form:"limit" json:"limit" query:"limit"`
	Page   int    `form:"page" json:"page" query:"page"`
	Sort   string `form:"sort" json:"sort" query:"sort"`
	Search string `form:"search" json:"search" query:"search"`
}
