package response

type Result struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// GostTermResult - result data structure for post term endpoint
type PostTermResult struct {
	Result
}

// GetTermResult - result data structure for get term endpoint
type GetTermResult struct {
	Payload interface{} `json:"payload,omitempty"`
	Result
}
