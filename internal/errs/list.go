package errs

import "errors"

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInvalidFieldValue = errors.New("invalid field value")
	ErrNotFound = errors.New("not found")
	ErrUserNotFound = errors.New("user not found")
	ErrVideoNotFound = errors.New("video not found")
	ErrCommentNotFound = errors.New("comment not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrChatNotFound = errors.New("chat not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrIncorrectUsernameOrPassword = errors.New("incorrect username or password")
	ErrInvalidToken = errors.New("invalid token")
	ErrAccessDenied = errors.New("access denied")
	ErrCannotSubscribeToYourself = errors.New("cannot subscribe to yourself")
	ErrCannotDonateToYourself = errors.New("cannot donate to yourself")
	ErrCannotCreateChatWithYourself = errors.New("cannot create chat with yourself")
	ErrPlaylistNotFound = errors.New("playlist not found")
)
