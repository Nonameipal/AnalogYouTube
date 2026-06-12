package dto

type CreateDonationRequest struct {
	ReceiverID int `json:"receiver_id"`
	VideoID *int `json:"video_id"`
	Amount float64 `json:"amount"`
	Message string `json:"message"`
}
