package transfer

type VideoTransfer struct {
	Category    string `json:"category"`
	Description string `json:"description"`
}
type VideoResponseTransfer struct {
	VideoID string `json:"video_id"`
}
