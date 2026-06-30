package repository

import (
	"context"

	dbModels "github.com/Nonameipal/AnalogYouTube/internal/models/db"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (r *Repository) CreateDonation(donation domain.Donation) (domain.Donation, error) {
	ctx := context.Background()
	var dbDonation dbModels.Donation
	err := r.db.QueryRow(ctx,
		`INSERT INTO donations (sender_id, receiver_id, video_id, amount, message)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''))
		RETURNING id, sender_id, receiver_id, video_id, amount, message, created_at`,
		donation.SenderID,
		donation.ReceiverID,
		donation.VideoID,
		donation.Amount,
		donation.Message,
	).Scan(&dbDonation.ID, &dbDonation.SenderID, &dbDonation.ReceiverID, &dbDonation.VideoID,
		&dbDonation.Amount, &dbDonation.Message, &dbDonation.CreatedAt)
	if err != nil {
		return domain.Donation{}, r.translateError(err)
	}

	return dbDonation.ToDomain(), nil
}

func (r *Repository) GetSentDonations(senderID int) ([]domain.Donation, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, sender_id, receiver_id, video_id, amount, message, created_at
		FROM donations
		WHERE sender_id = $1
		ORDER BY created_at DESC`, senderID)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbDonations []dbModels.Donation
	for rows.Next() {
		var donation dbModels.Donation
		if err := rows.Scan(&donation.ID, &donation.SenderID, &donation.ReceiverID, &donation.VideoID,
			&donation.Amount, &donation.Message, &donation.CreatedAt); err != nil {
			return nil, r.translateError(err)
		}
		dbDonations = append(dbDonations, donation)
	}

	if err := rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return donationsToDomain(dbDonations), nil
}

func (r *Repository) GetReceivedDonations(receiverID int) ([]domain.Donation, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, sender_id, receiver_id, video_id, amount, message, created_at
		FROM donations
		WHERE receiver_id = $1
		ORDER BY created_at DESC`, receiverID)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbDonations []dbModels.Donation
	for rows.Next() {
		var donation dbModels.Donation
		if err := rows.Scan(&donation.ID, &donation.SenderID, &donation.ReceiverID, &donation.VideoID,
			&donation.Amount, &donation.Message, &donation.CreatedAt); err != nil {
			return nil, r.translateError(err)
		}
		dbDonations = append(dbDonations, donation)
	}

	if err := rows.Err(); err != nil {
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
