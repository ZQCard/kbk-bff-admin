package paginationhelper

func Pagination(params map[string]interface{}) (page int64, page_size int64) {
	pageParams := params["Page"]
	page = 1
	if pageParams == nil || pageParams.(int64) == 0 {
	} else {
		page = pageParams.(int64)
	}
	page_sizeParams := params["PageSize"]
	page_size = 10
	if page_sizeParams == nil || page_sizeParams.(int64) == 0 {
	} else {
		page_size = page_sizeParams.(int64)
	}
	return page, page_size
}
