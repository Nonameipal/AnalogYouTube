package usecase

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (uc *Usecase) SendChatRequest(senderID int, receiverID int) (domain.ChatRequest, error) {
	if senderID <= 0 || receiverID <= 0 {
		return domain.ChatRequest{}, errs.ErrInvalidFieldValue
	}
	if senderID == receiverID {
		return domain.ChatRequest{}, errs.ErrCannotCreateChatWithYourself
	}

	if _, err := uc.GetUserByID(senderID); err != nil {
		return domain.ChatRequest{}, err
	}
	if _, err := uc.GetUserByID(receiverID); err != nil {
		return domain.ChatRequest{}, err
	}

	if _, err := uc.repository.GetChatBetweenUsers(senderID, receiverID); err == nil {
		return domain.ChatRequest{}, errs.ErrChatAlreadyExists
	} else if !errors.Is(err, errs.ErrNotFound) {
		return domain.ChatRequest{}, err
	}

	return uc.repository.CreateChatRequest(domain.ChatRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     domain.ChatRequestStatusPending,
	})
}

func (uc *Usecase) GetIncomingChatRequests(userID int) ([]domain.ChatRequest, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetUserByID(userID); err != nil {
		return nil, err
	}

	return uc.repository.GetIncomingChatRequests(userID)
}

func (uc *Usecase) GetOutgoingChatRequests(userID int) ([]domain.ChatRequest, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetUserByID(userID); err != nil {
		return nil, err
	}

	return uc.repository.GetOutgoingChatRequests(userID)
}

func (uc *Usecase) AcceptChatRequest(userID int, requestID int) (domain.AcceptedChatRequest, error) {
	request, err := uc.getPendingChatRequestForReceiver(userID, requestID)
	if err != nil {
		return domain.AcceptedChatRequest{}, err
	}

	acceptedRequest, chat, err := uc.repository.AcceptChatRequest(request.ID, request.SenderID, request.ReceiverID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.AcceptedChatRequest{}, errs.ErrChatRequestAlreadyAnswered
		}
		return domain.AcceptedChatRequest{}, err
	}

	return domain.AcceptedChatRequest{
		Request: acceptedRequest,
		Chat:    chat,
	}, nil
}

func (uc *Usecase) RejectChatRequest(userID int, requestID int) (domain.ChatRequest, error) {
	request, err := uc.getPendingChatRequestForReceiver(userID, requestID)
	if err != nil {
		return domain.ChatRequest{}, err
	}

	rejectedRequest, err := uc.repository.RejectChatRequest(request.ID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.ChatRequest{}, errs.ErrChatRequestAlreadyAnswered
		}
		return domain.ChatRequest{}, err
	}

	return rejectedRequest, nil
}

func (uc *Usecase) GetUserChats(userID int) ([]domain.Chat, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetUserByID(userID); err != nil {
		return nil, err
	}

	return uc.repository.GetUserChats(userID)
}

func (uc *Usecase) GetChatMessages(userID int, chatID int) ([]domain.ChatMessage, error) {
	if err := uc.EnsureUserCanAccessChat(userID, chatID); err != nil {
		return nil, err
	}

	return uc.repository.GetChatMessages(chatID)
}

func (uc *Usecase) CreateChatMessage(userID int, chatID int, message domain.ChatMessage) (domain.ChatMessage, error) {
	message.Text = strings.TrimSpace(message.Text)
	if userID <= 0 || chatID <= 0 || message.Text == "" {
		return domain.ChatMessage{}, errs.ErrInvalidFieldValue
	}

	if err := uc.EnsureUserCanAccessChat(userID, chatID); err != nil {
		return domain.ChatMessage{}, err
	}

	message.ChatID = chatID
	message.SenderID = userID

	return uc.repository.CreateChatMessage(message)
}

func (uc *Usecase) EnsureUserCanAccessChat(userID int, chatID int) error {
	if userID <= 0 || chatID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	chat, err := uc.repository.GetChatByID(chatID)
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

func (uc *Usecase) getPendingChatRequestForReceiver(userID int, requestID int) (domain.ChatRequest, error) {
	if userID <= 0 || requestID <= 0 {
		return domain.ChatRequest{}, errs.ErrInvalidFieldValue
	}

	request, err := uc.repository.GetChatRequestByID(requestID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.ChatRequest{}, errs.ErrChatRequestNotFound
		}
		return domain.ChatRequest{}, err
	}

	if request.ReceiverID != userID {
		return domain.ChatRequest{}, errs.ErrAccessDenied
	}
	if request.Status != domain.ChatRequestStatusPending {
		return domain.ChatRequest{}, errs.ErrChatRequestAlreadyAnswered
	}

	return request, nil
}
