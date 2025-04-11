package http

type CreateRequest struct {
	Type  string `json:"type"`
	PVZID string `json:"pvzId"`
}