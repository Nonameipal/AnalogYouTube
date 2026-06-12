package repository

import (
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/models/db"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (r *Repository) CreateDonation(donation domain.Donation) (domain.Donation, error) {
	var dbDonation dbModels.Donation
	err := r.db.Get(&dbDonation, 
		`INSERT INTO donations (sender_id, receiver_id, video_id, amount, message)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''))
		RETURNING id, sender_id, receiver_id, video_id, amount, message, created_at`,
		donation.SenderID,
		donation.ReceiverID,
		donation.VideoID,
		donation.Amount,
		donation.Message,
	)
	if err != nil {
		return domain.Donation{}, r.translateError(err)
	}

	return dbDonation.ToDomain(), nil
}

func (r *Repository) GetSentDonations(senderID int) ([]domain.Donation, error) {
	var dbDonations []dbModels.Donation
	err := r.db.Select(&dbDonations, 
		`SELECT id, sender_id, receiver_id, video_id, amount, message, created_at
		FROM donations
		WHERE sender_id = $1
		ORDER BY created_at DESC`, senderID)
	if err != nil {
		return nil, r.translateError(err)
	}

	return donationsToDomain(dbDonations), nil
}

func (r *Repository) GetReceivedDonations(receiverID int) ([]domain.Donation, error) {
	var dbDonations []dbModels.Donation
	err := r.db.Select(&dbDonations, 
		`SELECT id, sender_id, receiver_id, video_id, amount, message, created_at
		FROM donations
		WHERE receiver_id = $1
		ORDER BY created_at DESC`, receiverID)
	if err != nil {
		return nil, r.translateError(err)
	}

	return donationsToDomain(dbDonations), nil
}

func (r *Repository) GetUserDonations(userID int) ([]domain.Donation, error) {
	return r.GetReceivedDonations(userID)
}

func donationsToDomain(dbDonations []dbModels.Donation) []domain.Donation {
	donations := make([]domain.Donation, 0, len(dbDonations))
	for _, donation := range dbDonations {
		donations = append(donations, donation.ToDomain())
	}

	return donations
}
