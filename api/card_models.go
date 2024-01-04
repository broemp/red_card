package api

type createCardRequest struct {
	Color       string `json:"color" binding:"required,oneof=red yellow blue"`
	Accused     int64  `json:"accused" binding:"required"`
	Description string `json:"description"`
	Event       int64  `json:"event"`
}

type listCardRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}
