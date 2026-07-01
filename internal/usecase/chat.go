package usecase

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (s *Service) CreateOrGetChat(userID int, secondUserID int) (domain.Chat, error) {
	if userID <= 0 || secondUserID <= 0 {
		return domain.Chat{}, errs.ErrInvalidFieldValue
	}
	if userID == secondUserID {
		return domain.Chat{}, errs.ErrCannotCreateChatWithYourself
	}

	if _, err := s.GetUserByID(secondUserID); err != nil {
		return domain.Chat{}, err
	}

	return s.repository.CreateOrGetChat(userID, secondUserID)
}

func (s *Service) GetUserChats(userID int) ([]domain.Chat, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	return s.repository.GetUserChats(userID)
}

func (s *Service) GetChatMessages(userID int, chatID int) ([]domain.ChatMessage, error) {
	if err := s.EnsureUserCanAccessChat(userID, chatID); err != nil {
		return nil, err
	}

	return s.repository.GetChatMessages(chatID)
}

func (s *Service) CreateChatMessage(userID int, chatID int, message domain.ChatMessage) (domain.ChatMessage, error) {
	message.Text = strings.TrimSpace(message.Text)
	if userID <= 0 || chatID <= 0 || message.Text == "" {
		return domain.ChatMessage{}, errs.ErrInvalidFieldValue
	}

	if err := s.EnsureUserCanAccessChat(userID, chatID); err != nil {
		return domain.ChatMessage{}, err
	}

	message.ChatID = chatID
	message.SenderID = userID

	return s.repository.CreateChatMessage(message)
}

func (s *Service) EnsureUserCanAccessChat(userID int, chatID int) error {
	if userID <= 0 || chatID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	chat, err := s.repository.GetChatByID(chatID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrChatNotFound
		}
		return err
	}

	if chat.FirstUserID != userID && chat.SecondUserID != userID {
		return errs.ErrAccessDenied
	}

	return nil
}
