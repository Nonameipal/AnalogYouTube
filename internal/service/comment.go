package service

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (s *Service) CreateComment(userID int, videoID int, comment domain.Comment) (domain.Comment, error) {
	comment.Text = strings.TrimSpace(comment.Text)

	if userID <= 0 || videoID <= 0 || comment.Text == "" {
		return domain.Comment{}, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetVideoByID(videoID); err != nil {
		return domain.Comment{}, err
	}

	comment.UserID = userID
	comment.VideoID = videoID
	comment.Status = domain.CommentStatusActive

	return s.repository.CreateComment(comment)
}

func (s *Service) GetVideoComments(videoID int) ([]domain.Comment, error) {
	if videoID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetVideoByID(videoID); err != nil {
		return nil, err
	}

	return s.repository.GetVideoComments(videoID)
}

func (s *Service) UpdateComment(userID int, userRole string, comment domain.Comment) (domain.Comment, error) {
	comment.Text = strings.TrimSpace(comment.Text)

	if userID <= 0 || comment.ID <= 0 || comment.Text == "" {
		return domain.Comment{}, errs.ErrInvalidFieldValue
	}

	oldComment, err := s.repository.GetCommentByID(comment.ID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Comment{}, errs.ErrCommentNotFound
		}
		return domain.Comment{}, err
	}

	if oldComment.UserID != userID && userRole != domain.AdminRole {
		return domain.Comment{}, errs.ErrAccessDenied
	}

	return s.repository.UpdateComment(comment)
}

func (s *Service) DeleteComment(userID int, userRole string, commentID int) error {
	if userID <= 0 || commentID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	comment, err := s.repository.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrCommentNotFound
		}
		return err
	}

	if comment.UserID != userID && userRole != domain.AdminRole {
		return errs.ErrAccessDenied
	}

	if err = s.repository.DeleteComment(commentID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrCommentNotFound
		}
		return err
	}

	return nil
}
