package twitter

type (
	Tweet struct {
		Id   string `json:"id"`
		Text string `json:"text"`
	}

	GetTweetsResponse struct {
		Data []Tweet `json:"data"`
	}
)
