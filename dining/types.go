package dining

type RequestItem struct {
	ItemID   int    `json:"item_id"`
	Quantity int    `json:"quantity"`
	Note     string `json:"notes"`
}
