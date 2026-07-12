package usecase

import (
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/infrastructure/ffmpeg"
)

type fakeRepository struct {
	usersByID       map[int]domain.User
	usersByUsername map[string]domain.User
	usersByEmail    map[string]domain.User

	videosByID     map[int]domain.Video
	videoQualities map[int][]domain.VideoQuality
	categories     map[string]domain.Category
	tags           []domain.Tag

	commentsByID   map[int]domain.Comment
	chatsByID      map[int]domain.Chat
	chatRequests   map[int]domain.ChatRequest
	playlistsByID  map[int]domain.Playlist
	playlistVideos map[int][]domain.Video

	recommendedVideos []domain.Video
	searchVideos      []domain.Video
	subscribers       []domain.User
	subscriptions     []domain.User

	createdUser        domain.User
	updatedUser        domain.User
	createdVideo       domain.Video
	updatedVideo       domain.Video
	createdCategory    domain.Category
	updatedCategory    domain.Category
	createdDonation    domain.Donation
	createdComment     domain.Comment
	updatedComment     domain.Comment
	createdChatRequest domain.ChatRequest
	acceptedRequest    domain.ChatRequest
	acceptedChat       domain.Chat
	rejectedRequest    domain.ChatRequest
	createdMessage     domain.ChatMessage
	createdPlaylist    domain.Playlist
	updatedPlaylist    domain.Playlist
	createdQuality     domain.VideoQuality

	deletedVideoID         int
	archivedVideoID        int
	archiveURL             string
	deletedCategoryID      int
	likedUserID            int
	likedVideoID           int
	unlikedUserID          int
	unlikedVideoID         int
	subscribedUserID       int
	subscribedAuthorID     int
	unsubscribedUserID     int
	unsubscribedAuthorID   int
	deletedCommentID       int
	deletedPlaylistID      int
	addedPlaylistID        int
	addedVideoID           int
	removedPlaylistID      int
	removedVideoID         int
	incrementedVideoID     int
	videoLikesCount        int
	videoLikedByUser       bool
	subscribersCount       int
	subscriptionsCount     int
	isSubscribed           bool
	nextID                 int
	createUserErr          error
	createVideoErr         error
	updateVideoErr         error
	deleteVideoErr         error
	updateCategoryErr      error
	deleteCategoryErr      error
	createDonationErr      error
	createCommentErr       error
	updateCommentErr       error
	deleteCommentErr       error
	createChatRequestErr   error
	acceptChatRequestErr   error
	rejectChatRequestErr   error
	createChatMessageErr   error
	createPlaylistErr      error
	updatePlaylistErr      error
	deletePlaylistErr      error
	addVideoPlaylistErr    error
	removeVideoPlaylistErr error
}

func newTestUsecase(repo *fakeRepository) *Usecase {
	return NewUsecase(repo, ffmpeg.NewFFmpegSettings("ffmpeg"))
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		usersByID:          make(map[int]domain.User),
		usersByUsername:    make(map[string]domain.User),
		usersByEmail:       make(map[string]domain.User),
		videosByID:         make(map[int]domain.Video),
		videoQualities:     make(map[int][]domain.VideoQuality),
		categories:         make(map[string]domain.Category),
		tags:               make([]domain.Tag, 0),
		commentsByID:       make(map[int]domain.Comment),
		chatsByID:          make(map[int]domain.Chat),
		chatRequests:       make(map[int]domain.ChatRequest),
		playlistsByID:      make(map[int]domain.Playlist),
		playlistVideos:     make(map[int][]domain.Video),
		nextID:             1,
		videoLikesCount:    3,
		subscribersCount:   4,
		subscriptionsCount: 2,
	}
}

func (r *fakeRepository) CreateUser(user domain.User) error {
	if r.createUserErr != nil {
		return r.createUserErr
	}
	if user.ID == 0 {
		user.ID = r.nextID
		r.nextID++
	}
	r.createdUser = user
	r.usersByID[user.ID] = user
	r.usersByUsername[user.Username] = user
	if user.Email != "" {
		r.usersByEmail[user.Email] = user
	}
	return nil
}

func (r *fakeRepository) GetUserByUsername(username string) (domain.User, error) {
	user, ok := r.usersByUsername[username]
	if !ok {
		return domain.User{}, errs.ErrNotFound
	}
	return user, nil
}

func (r *fakeRepository) GetUserByEmail(email string) (domain.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return domain.User{}, errs.ErrNotFound
	}
	return user, nil
}

func (r *fakeRepository) GetUserByID(id int) (domain.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return domain.User{}, errs.ErrNotFound
	}
	return user, nil
}

func (r *fakeRepository) UpdateUserProfile(user domain.User) (domain.User, error) {
	r.updatedUser = user
	r.usersByID[user.ID] = user
	r.usersByUsername[user.Username] = user
	if user.Email != "" {
		r.usersByEmail[user.Email] = user
	}
	return user, nil
}

func (r *fakeRepository) CreateVideo(video domain.Video) (domain.Video, error) {
	if r.createVideoErr != nil {
		return domain.Video{}, r.createVideoErr
	}
	if video.ID == 0 {
		video.ID = r.nextID
		r.nextID++
	}
	r.createdVideo = video
	r.videosByID[video.ID] = video
	return video, nil
}

func (r *fakeRepository) GetAllVideos() ([]domain.Video, error) {
	videos := make([]domain.Video, 0, len(r.videosByID))
	for _, video := range r.videosByID {
		videos = append(videos, video)
	}
	return videos, nil
}

func (r *fakeRepository) GetRecommendedVideos() ([]domain.Video, error) {
	return r.recommendedVideos, nil
}

func (r *fakeRepository) GetVideoByID(id int) (domain.Video, error) {
	video, ok := r.videosByID[id]
	if !ok {
		return domain.Video{}, errs.ErrNotFound
	}
	return video, nil
}

func (r *fakeRepository) GetVideosByAuthorID(authorID int) ([]domain.Video, error) {
	var videos []domain.Video
	for _, video := range r.videosByID {
		if video.AuthorID == authorID {
			videos = append(videos, video)
		}
	}
	return videos, nil
}

func (r *fakeRepository) SearchVideosByTitle(title string) ([]domain.Video, error) {
	return r.searchVideos, nil
}

func (r *fakeRepository) IncrementVideoViews(id int) error {
	if _, ok := r.videosByID[id]; !ok {
		return errs.ErrNotFound
	}
	r.incrementedVideoID = id
	return nil
}

func (r *fakeRepository) UpdateVideo(video domain.Video) (domain.Video, error) {
	if r.updateVideoErr != nil {
		return domain.Video{}, r.updateVideoErr
	}
	r.updatedVideo = video
	r.videosByID[video.ID] = video
	return video, nil
}

func (r *fakeRepository) CreateVideoQuality(quality domain.VideoQuality) (domain.VideoQuality, error) {
	r.createdQuality = quality
	r.videoQualities[quality.VideoID] = append(r.videoQualities[quality.VideoID], quality)
	return quality, nil
}

func (r *fakeRepository) GetVideoQualities(videoID int) ([]domain.VideoQuality, error) {
	return r.videoQualities[videoID], nil
}

func (r *fakeRepository) DeleteVideo(id int) error {
	if r.deleteVideoErr != nil {
		return r.deleteVideoErr
	}
	if _, ok := r.videosByID[id]; !ok {
		return errs.ErrNotFound
	}
	r.deletedVideoID = id
	delete(r.videosByID, id)
	return nil
}

func (r *fakeRepository) ArchiveDeletedVideo(videoID int, archiveURL string) error {
	if _, ok := r.videosByID[videoID]; !ok {
		return errs.ErrNotFound
	}
	r.archivedVideoID = videoID
	r.archiveURL = archiveURL
	return nil
}

func (r *fakeRepository) CreateCategory(category domain.Category) (domain.Category, error) {
	if category.ID == 0 {
		category.ID = r.nextID
		r.nextID++
	}
	r.createdCategory = category
	r.categories[category.Name] = category
	return category, nil
}

func (r *fakeRepository) GetAllCategories() ([]domain.Category, error) {
	categories := make([]domain.Category, 0, len(r.categories))
	for _, category := range r.categories {
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *fakeRepository) GetCategoryByName(name string) (domain.Category, error) {
	category, ok := r.categories[name]
	if !ok {
		return domain.Category{}, errs.ErrNotFound
	}
	return category, nil
}

func (r *fakeRepository) UpdateCategory(category domain.Category) (domain.Category, error) {
	if r.updateCategoryErr != nil {
		return domain.Category{}, r.updateCategoryErr
	}
	r.updatedCategory = category
	r.categories[category.Name] = category
	return category, nil
}

func (r *fakeRepository) DeleteCategory(id int) error {
	if r.deleteCategoryErr != nil {
		return r.deleteCategoryErr
	}
	r.deletedCategoryID = id
	return nil
}

func (r *fakeRepository) CreateDonation(donation domain.Donation) (domain.Donation, error) {
	if r.createDonationErr != nil {
		return domain.Donation{}, r.createDonationErr
	}
	if donation.ID == 0 {
		donation.ID = r.nextID
		r.nextID++
	}
	r.createdDonation = donation
	return donation, nil
}

func (r *fakeRepository) GetSentDonations(senderID int) ([]domain.Donation, error) {
	return []domain.Donation{{SenderID: senderID}}, nil
}

func (r *fakeRepository) GetReceivedDonations(receiverID int) ([]domain.Donation, error) {
	return []domain.Donation{{ReceiverID: receiverID}}, nil
}

func (r *fakeRepository) GetUserDonations(userID int) ([]domain.Donation, error) {
	return []domain.Donation{{ReceiverID: userID}}, nil
}

func (r *fakeRepository) LikeVideo(userID int, videoID int) error {
	r.likedUserID = userID
	r.likedVideoID = videoID
	return nil
}

func (r *fakeRepository) UnlikeVideo(userID int, videoID int) error {
	r.unlikedUserID = userID
	r.unlikedVideoID = videoID
	return nil
}

func (r *fakeRepository) GetVideoLikesCount(videoID int) (int, error) {
	return r.videoLikesCount, nil
}

func (r *fakeRepository) IsVideoLikedByUser(userID int, videoID int) (bool, error) {
	return r.videoLikedByUser, nil
}

func (r *fakeRepository) SubscribeToUser(subscriberID int, authorID int) error {
	r.subscribedUserID = subscriberID
	r.subscribedAuthorID = authorID
	return nil
}

func (r *fakeRepository) UnsubscribeFromUser(subscriberID int, authorID int) error {
	r.unsubscribedUserID = subscriberID
	r.unsubscribedAuthorID = authorID
	return nil
}

func (r *fakeRepository) GetSubscribersCount(authorID int) (int, error) {
	return r.subscribersCount, nil
}

func (r *fakeRepository) GetSubscriptionsCount(subscriberID int) (int, error) {
	return r.subscriptionsCount, nil
}

func (r *fakeRepository) IsSubscribed(subscriberID int, authorID int) (bool, error) {
	return r.isSubscribed, nil
}

func (r *fakeRepository) GetSubscribers(authorID int) ([]domain.User, error) {
	return r.subscribers, nil
}

func (r *fakeRepository) GetSubscriptions(subscriberID int) ([]domain.User, error) {
	return r.subscriptions, nil
}

func (r *fakeRepository) CreateComment(comment domain.Comment) (domain.Comment, error) {
	if r.createCommentErr != nil {
		return domain.Comment{}, r.createCommentErr
	}
	if comment.ID == 0 {
		comment.ID = r.nextID
		r.nextID++
	}
	r.createdComment = comment
	r.commentsByID[comment.ID] = comment
	return comment, nil
}

func (r *fakeRepository) GetVideoComments(videoID int) ([]domain.Comment, error) {
	var comments []domain.Comment
	for _, comment := range r.commentsByID {
		if comment.VideoID == videoID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

func (r *fakeRepository) GetCommentByID(commentID int) (domain.Comment, error) {
	comment, ok := r.commentsByID[commentID]
	if !ok {
		return domain.Comment{}, errs.ErrNotFound
	}
	return comment, nil
}

func (r *fakeRepository) UpdateComment(comment domain.Comment) (domain.Comment, error) {
	if r.updateCommentErr != nil {
		return domain.Comment{}, r.updateCommentErr
	}
	r.updatedComment = comment
	r.commentsByID[comment.ID] = comment
	return comment, nil
}

func (r *fakeRepository) DeleteComment(commentID int) error {
	if r.deleteCommentErr != nil {
		return r.deleteCommentErr
	}
	if _, ok := r.commentsByID[commentID]; !ok {
		return errs.ErrNotFound
	}
	r.deletedCommentID = commentID
	delete(r.commentsByID, commentID)
	return nil
}

func (r *fakeRepository) GetChatByID(chatID int) (domain.Chat, error) {
	chat, ok := r.chatsByID[chatID]
	if !ok {
		return domain.Chat{}, errs.ErrNotFound
	}
	return chat, nil
}

func (r *fakeRepository) GetUserChats(userID int) ([]domain.Chat, error) {
	var chats []domain.Chat
	for _, chat := range r.chatsByID {
		if chat.FirstUserID == userID || chat.SecondUserID == userID {
			chats = append(chats, chat)
		}
	}
	return chats, nil
}

func (r *fakeRepository) GetChatBetweenUsers(firstUserID int, secondUserID int) (domain.Chat, error) {
	for _, chat := range r.chatsByID {
		if (chat.FirstUserID == firstUserID && chat.SecondUserID == secondUserID) ||
			(chat.FirstUserID == secondUserID && chat.SecondUserID == firstUserID) {
			return chat, nil
		}
	}
	return domain.Chat{}, errs.ErrNotFound
}

func (r *fakeRepository) CreateChatRequest(request domain.ChatRequest) (domain.ChatRequest, error) {
	if r.createChatRequestErr != nil {
		return domain.ChatRequest{}, r.createChatRequestErr
	}
	if request.ID == 0 {
		request.ID = r.nextID
		r.nextID++
	}
	r.createdChatRequest = request
	r.chatRequests[request.ID] = request
	return request, nil
}

func (r *fakeRepository) GetChatRequestByID(requestID int) (domain.ChatRequest, error) {
	request, ok := r.chatRequests[requestID]
	if !ok {
		return domain.ChatRequest{}, errs.ErrNotFound
	}
	return request, nil
}

func (r *fakeRepository) GetIncomingChatRequests(userID int) ([]domain.ChatRequest, error) {
	var requests []domain.ChatRequest
	for _, request := range r.chatRequests {
		if request.ReceiverID == userID {
			requests = append(requests, request)
		}
	}
	return requests, nil
}

func (r *fakeRepository) GetOutgoingChatRequests(userID int) ([]domain.ChatRequest, error) {
	var requests []domain.ChatRequest
	for _, request := range r.chatRequests {
		if request.SenderID == userID {
			requests = append(requests, request)
		}
	}
	return requests, nil
}

func (r *fakeRepository) AcceptChatRequest(requestID int, firstUserID int, secondUserID int) (domain.ChatRequest, domain.Chat, error) {
	if r.acceptChatRequestErr != nil {
		return domain.ChatRequest{}, domain.Chat{}, r.acceptChatRequestErr
	}
	request := r.chatRequests[requestID]
	request.Status = domain.ChatRequestStatusAccepted
	chat := domain.Chat{ID: r.nextID, FirstUserID: firstUserID, SecondUserID: secondUserID}
	r.nextID++
	r.acceptedRequest = request
	r.acceptedChat = chat
	r.chatRequests[requestID] = request
	r.chatsByID[chat.ID] = chat
	return request, chat, nil
}

func (r *fakeRepository) RejectChatRequest(requestID int) (domain.ChatRequest, error) {
	if r.rejectChatRequestErr != nil {
		return domain.ChatRequest{}, r.rejectChatRequestErr
	}
	request := r.chatRequests[requestID]
	request.Status = domain.ChatRequestStatusRejected
	r.rejectedRequest = request
	r.chatRequests[requestID] = request
	return request, nil
}

func (r *fakeRepository) CreateChatMessage(message domain.ChatMessage) (domain.ChatMessage, error) {
	if r.createChatMessageErr != nil {
		return domain.ChatMessage{}, r.createChatMessageErr
	}
	if message.ID == 0 {
		message.ID = r.nextID
		r.nextID++
	}
	r.createdMessage = message
	return message, nil
}

func (r *fakeRepository) GetChatMessages(chatID int) ([]domain.ChatMessage, error) {
	return []domain.ChatMessage{{ChatID: chatID, Text: "hello"}}, nil
}

func (r *fakeRepository) CreatePlaylist(playlist domain.Playlist) (domain.Playlist, error) {
	if r.createPlaylistErr != nil {
		return domain.Playlist{}, r.createPlaylistErr
	}
	if playlist.ID == 0 {
		playlist.ID = r.nextID
		r.nextID++
	}
	r.createdPlaylist = playlist
	r.playlistsByID[playlist.ID] = playlist
	return playlist, nil
}

func (r *fakeRepository) GetPlaylistByID(id int) (domain.Playlist, error) {
	playlist, ok := r.playlistsByID[id]
	if !ok {
		return domain.Playlist{}, errs.ErrNotFound
	}
	return playlist, nil
}

func (r *fakeRepository) GetUserPlaylists(userID int) ([]domain.Playlist, error) {
	var playlists []domain.Playlist
	for _, playlist := range r.playlistsByID {
		if playlist.UserID == userID {
			playlists = append(playlists, playlist)
		}
	}
	return playlists, nil
}

func (r *fakeRepository) UpdatePlaylist(playlist domain.Playlist) (domain.Playlist, error) {
	if r.updatePlaylistErr != nil {
		return domain.Playlist{}, r.updatePlaylistErr
	}
	r.updatedPlaylist = playlist
	r.playlistsByID[playlist.ID] = playlist
	return playlist, nil
}

func (r *fakeRepository) DeletePlaylist(id int) error {
	if r.deletePlaylistErr != nil {
		return r.deletePlaylistErr
	}
	if _, ok := r.playlistsByID[id]; !ok {
		return errs.ErrNotFound
	}
	r.deletedPlaylistID = id
	delete(r.playlistsByID, id)
	return nil
}

func (r *fakeRepository) AddVideoToPlaylist(playlistID int, videoID int) error {
	if r.addVideoPlaylistErr != nil {
		return r.addVideoPlaylistErr
	}
	r.addedPlaylistID = playlistID
	r.addedVideoID = videoID
	return nil
}

func (r *fakeRepository) RemoveVideoFromPlaylist(playlistID int, videoID int) error {
	if r.removeVideoPlaylistErr != nil {
		return r.removeVideoPlaylistErr
	}
	r.removedPlaylistID = playlistID
	r.removedVideoID = videoID
	return nil
}

func (r *fakeRepository) GetPlaylistVideos(playlistID int) ([]domain.Video, error) {
	return r.playlistVideos[playlistID], nil
}

func (r *fakeRepository) GetAllTags() ([]domain.Tag, error) {
	return r.tags, nil
}

func (r *fakeRepository) GetVideoTags(videoID int) ([]domain.Tag, error) {
	return []domain.Tag{}, nil
}

func (r *fakeRepository) GetTagsByVideoIDs(videoIDs []int) (map[int][]domain.Tag, error) {
	result := make(map[int][]domain.Tag)
	for _, videoID := range videoIDs {
		result[videoID] = []domain.Tag{}
	}
	return result, nil
}

func (r *fakeRepository) ReplaceVideoTags(videoID int, tagNames []string) ([]domain.Tag, error) {
	var result []domain.Tag
	for i, name := range tagNames {
		result = append(result, domain.Tag{ID: i + 1, Name: name})
	}
	return result, nil
}

func (r *fakeRepository) GetRecommendationCandidateVideos(limit int) ([]domain.Video, error) {
	videos := make([]domain.Video, 0, limit)
	count := 0
	for _, video := range r.videosByID {
		if count >= limit {
			break
		}
		videos = append(videos, video)
		count++
	}
	return videos, nil
}

func (r *fakeRepository) GetVideosByIDs(ids []int) ([]domain.Video, error) {
	var result []domain.Video
	for _, id := range ids {
		if video, ok := r.videosByID[id]; ok {
			result = append(result, video)
		}
	}
	return result, nil
}

func (r *fakeRepository) SaveVideoWatchProgress(progress domain.VideoWatchHistory) (domain.VideoWatchHistory, error) {
	return progress, nil
}

func (r *fakeRepository) GetUserWatchHistory(userID int, limit int) ([]domain.VideoWatchHistory, error) {
	return []domain.VideoWatchHistory{}, nil
}

func (r *fakeRepository) SaveVideoSearchHistory(history domain.VideoSearchHistory) error {
	return nil
}

func (r *fakeRepository) GetUserSearchHistory(userID int, limit int) ([]domain.VideoSearchHistory, error) {
	return []domain.VideoSearchHistory{}, nil
}

func (r *fakeRepository) GetUserLikedVideoIDs(userID int) ([]int, error) {
	return []int{}, nil
}

func (r *fakeRepository) GetSubscribedAuthorIDs(subscriberID int) ([]int, error) {
	return []int{}, nil
}

func (r *fakeRepository) GetUserCommentedVideoIDs(userID int) ([]int, error) {
	return []int{}, nil
}

func (r *fakeRepository) GetUserPlaylistVideoIDs(userID int) ([]int, error) {
	return []int{}, nil
}

func (r *fakeRepository) GetVideoPopularity(videoIDs []int) (map[int]domain.VideoPopularity, error) {
	return make(map[int]domain.VideoPopularity), nil
}
