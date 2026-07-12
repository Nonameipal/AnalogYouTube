package usecase

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (uc *Usecase) CreateComment(userID int, videoID int, comment domain.Comment) (domain.Comment, error) {
	comment.Text = strings.TrimSpace(comment.Text)

	if userID <= 0 || videoID <= 0 || comment.Text == "" {
		return domain.Comment{}, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetVideoByID(videoID); err != nil {
		return domain.Comment{}, err
	}

	if comment.ParentID != nil {
		parent, err := uc.repository.GetCommentByID(*comment.ParentID)
		if err != nil {
			return domain.Comment{}, errs.ErrCommentNotFound
		}
		if parent.VideoID != videoID {
			return domain.Comment{}, errs.ErrInvalidFieldValue
		}
	}

	comment.UserID = userID
	comment.VideoID = videoID
	comment.Status = domain.CommentStatusActive

	return uc.repository.CreateComment(comment)
}

func (uc *Usecase) GetVideoComments(videoID int) ([]domain.Comment, error) {
	if videoID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetVideoByID(videoID); err != nil {
		return nil, err
	}

	return uc.repository.GetVideoComments(videoID)
}

func (uc *Usecase) UpdateComment(userID int, userRole string, comment domain.Comment) (domain.Comment, error) {
	comment.Text = strings.TrimSpace(comment.Text)

	if userID <= 0 || comment.ID <= 0 || comment.Text == "" {
		return domain.Comment{}, errs.ErrInvalidFieldValue
	}

	oldComment, err := uc.repository.GetCommentByID(comment.ID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Comment{}, errs.ErrCommentNotFound
		}
		return domain.Comment{}, err
	}

	if oldComment.UserID != userID && userRole != domain.AdminRole {
		return domain.Comment{}, errs.ErrAccessDenied
	}

	return uc.repository.UpdateComment(comment)
}

func (uc *Usecase) DeleteComment(userID int, userRole string, commentID int) error {
	if userID <= 0 || commentID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	comment, err := uc.repository.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrCommentNotFound
		}
		return err
	}

	if comment.UserID != userID && userRole != domain.AdminRole {
		return errs.ErrAccessDenied
	}

	if err = uc.repository.DeleteComment(commentID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrCommentNotFound
		}
		return err
	}

	return nil
}
