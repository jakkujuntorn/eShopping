package models

import ()

// ปั้น pagination ก่อนส่งเข้า db
type Pagination struct {
	Page   int
	// Total  int
	Limit  int
	Offset int // page-1 * limit
}

// ปั้น pagination ก่อนส่งเข้า db
func PaginationDB(page, limit, offset int) Pagination {
	return Pagination{
		Page:   page,
		// Total:  total,
		Limit:  limit,
		Offset: offset,
	}
}

type PageResponse struct {
	Success string
	Total    int
	Page     int
	PageSize int
	Data     interface{}
}

// ปั้ม pagination ก่อนส่งไป fotn end
func PaginationResponse(message string, total, page, pagesize int, data interface{}) PageResponse {
	return PageResponse{
		Success:  message,
		Total:    total,
		Page:     page,
		PageSize: pagesize,
		Data:     data,
	}
}
