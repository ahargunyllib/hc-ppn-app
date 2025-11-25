package dto

type PaginationResponse struct {
	TotalData int64 `json:"total_data"`
	TotalPage int   `json:"total_page"`
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
}

func NewPaginationResponse(totalData int64, page, limit int) PaginationResponse {
	totalPage := int(totalData) / limit
	if int(totalData)%limit != 0 {
		totalPage++
	}

	return PaginationResponse{
		TotalData: totalData,
		TotalPage: totalPage,
		Page:      page,
		Limit:     limit,
	}
}
