package http

type CrackHashRequest struct {
	Hash      string `json:"hash"`
	MaxLength int    `json:"maxLength"`
}

type CrackHashResponse struct {
	RequestId string `json:"requestId"`
}

type CrackHashStatus struct {
	Status   string   `json:"status"`
	Data     []string `json:"data"`
	Progress int      `json:"progress"`
}
