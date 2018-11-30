package model

type Item struct {
	ID string `json:"id"`
	Url string `json:"url"`
	Title string `json:"title"`
	LikesCount int `json:"likes_count"`
	UpdatedAt string `json:"updated_at"`
	
}

type ItemResponse struct {
	ID string `json:"id"`
	Url string `json:"url"`
	Title string `json:"title"`
	LikesCount int `json:"likes_count"`
	UpdatedAt string `json:"updated_at"`
}

func NewResponse(items []*Item) []*ItemResponse {
	var itemResponses []*ItemResponse
	for _,item := range items {
		itemResponse := &ItemResponse{
			ID:item.ID,
			Url:item.Url,
			Title: item.Title,
			LikesCount:item.LikesCount,
			UpdatedAt:item.UpdatedAt,
		}
		itemResponses = append(itemResponses,itemResponse)
	}

	return itemResponses
}