package dto

type CreateDonationRequest struct {
	ReceiverID  int     `json:"receiver_id"`
	RecipientID int     `json:"recipient_id"`
	VideoID     *int    `json:"video_id"`
	Amount      float64 `json:"amount"`
	Message     string  `json:"message"`
}

func (r CreateDonationRequest) TargetUserID() int {
	if r.ReceiverID > 0 {
		return r.ReceiverID
	}

	return r.RecipientID
}
