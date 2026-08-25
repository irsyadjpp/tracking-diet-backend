package usecase

import (
	"errors"
	"fmt"

	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
)

var (
	ErrConnectionNotFound  = errors.New("connection not found")
	ErrPostNotFound        = errors.New("post not found")
	ErrGroupNotFound       = errors.New("group not found")
	ErrAlreadyConnected    = errors.New("already connected")
	ErrAlreadyLiked       = errors.New("already liked")
)

// SocialUsecase handles social features business logic
type SocialUsecase struct {
	socialRepo domain.SocialRepository
	userRepo   domain.UserRepository
}

// NewSocialUsecase creates a new social use case
func NewSocialUsecase(socialRepo domain.SocialRepository, userRepo domain.UserRepository) *SocialUsecase {
	return &SocialUsecase{
		socialRepo: socialRepo,
		userRepo:   userRepo,
	}
}

// Connections
func (uc *SocialUsecase) SendConnectionRequest(userID, connectedID int, connectionType domain.ConnectionType) (*domain.Connection, error) {
	if userID == connectedID {
		return nil, errors.New("cannot connect to yourself")
	}

	// Check if connection already exists
	existing, err := uc.socialRepo.GetConnection(userID, connectedID)
	if err == nil && existing != nil {
		return nil, ErrAlreadyConnected
	}

	connection := &domain.Connection{
		UserID:      userID,
		ConnectedID: connectedID,
		Type:        connectionType,
		Status:      domain.ConnectionStatusPending,
	}

	if err := uc.socialRepo.CreateConnection(connection); err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	return connection, nil
}

func (uc *SocialUsecase) AcceptConnection(connectionID int) error {
	if err := uc.socialRepo.UpdateConnectionStatus(connectionID, domain.ConnectionStatusAccepted); err != nil {
		return fmt.Errorf("failed to accept connection: %w", err)
	}
	return nil
}

func (uc *SocialUsecase) RejectConnection(connectionID int) error {
	if err := uc.socialRepo.UpdateConnectionStatus(connectionID, domain.ConnectionStatusRejected); err != nil {
		return fmt.Errorf("failed to reject connection: %w", err)
	}
	return nil
}

func (uc *SocialUsecase) GetUserConnections(userID int) ([]domain.Connection, error) {
	connections, err := uc.socialRepo.GetUserConnections(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user connections: %w", err)
	}
	return connections, nil
}

func (uc *SocialUsecase) GetPendingConnections(userID int) ([]domain.Connection, error) {
	connections, err := uc.socialRepo.GetPendingConnections(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending connections: %w", err)
	}
	return connections, nil
}

func (uc *SocialUsecase) RemoveConnection(connectionID int) error {
	if err := uc.socialRepo.DeleteConnection(connectionID); err != nil {
		return fmt.Errorf("failed to remove connection: %w", err)
	}
	return nil
}

// Posts
func (uc *SocialUsecase) CreatePost(userID int, postType domain.PostType, title, content, imageURL string, isPublic bool) (*domain.Post, error) {
	post := &domain.Post{
		UserID:       userID,
		Type:         postType,
		Title:        title,
		Content:      content,
		ImageURL:     imageURL,
		IsPublic:     isPublic,
		LikesCount:   0,
		CommentsCount: 0,
	}

	if err := uc.socialRepo.CreatePost(post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (uc *SocialUsecase) GetPost(id int) (*domain.Post, error) {
	post, err := uc.socialRepo.GetPost(id)
	if err != nil {
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (uc *SocialUsecase) GetUserPosts(userID int) ([]domain.Post, error) {
	posts, err := uc.socialRepo.GetUserPosts(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}
	return posts, nil
}

func (uc *SocialUsecase) GetFeedPosts(userID int) ([]domain.Post, error) {
	posts, err := uc.socialRepo.GetFeedPosts(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed posts: %w", err)
	}
	return posts, nil
}

func (uc *SocialUsecase) UpdatePost(id int, title, content *string, imageURL *string, isPublic *bool) (*domain.Post, error) {
	post, err := uc.socialRepo.GetPost(id)
	if err != nil {
		return nil, ErrPostNotFound
	}

	if title != nil {
		post.Title = *title
	}
	if content != nil {
		post.Content = *content
	}
	if imageURL != nil {
		post.ImageURL = *imageURL
	}
	if isPublic != nil {
		post.IsPublic = *isPublic
	}

	if err := uc.socialRepo.UpdatePost(post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

func (uc *SocialUsecase) DeletePost(id int) error {
	if err := uc.socialRepo.DeletePost(id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

// Comments
func (uc *SocialUsecase) AddComment(postID, userID int, content string) (*domain.Comment, error) {
	comment := &domain.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}

	if err := uc.socialRepo.CreateComment(comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

func (uc *SocialUsecase) GetPostComments(postID int) ([]domain.Comment, error) {
	comments, err := uc.socialRepo.GetPostComments(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post comments: %w", err)
	}
	return comments, nil
}

func (uc *SocialUsecase) UpdateComment(id int, content string) (*domain.Comment, error) {
	comment := &domain.Comment{
		ID:      id,
		Content: content,
	}

	if err := uc.socialRepo.UpdateComment(comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return comment, nil
}

func (uc *SocialUsecase) DeleteComment(id int) error {
	if err := uc.socialRepo.DeleteComment(id); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

// Likes
func (uc *SocialUsecase) LikePost(postID, userID int) error {
	// Check if already liked
	alreadyLiked, err := uc.socialRepo.CheckUserLiked(postID, userID)
	if err != nil {
		return fmt.Errorf("failed to check like status: %w", err)
	}
	if alreadyLiked {
		return ErrAlreadyLiked
	}

	if err := uc.socialRepo.CreateLike(postID, userID); err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	return nil
}

func (uc *SocialUsecase) UnlikePost(postID, userID int) error {
	if err := uc.socialRepo.DeleteLike(postID, userID); err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}
	return nil
}

func (uc *SocialUsecase) GetPostLikes(postID int) ([]domain.Like, error) {
	likes, err := uc.socialRepo.GetPostLikes(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post likes: %w", err)
	}
	return likes, nil
}

// Groups
func (uc *SocialUsecase) CreateGroup(createdBy int, name, description, imageURL string, isPublic bool) (*domain.Group, error) {
	group := &domain.Group{
		Name:        name,
		Description: description,
		ImageURL:    imageURL,
		IsPublic:    isPublic,
		MemberCount: 1,
		CreatedBy:   createdBy,
	}

	if err := uc.socialRepo.CreateGroup(group); err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// Add creator as admin
	if err := uc.socialRepo.JoinGroup(group.ID, createdBy, "admin"); err != nil {
		return nil, fmt.Errorf("failed to add creator to group: %w", err)
	}

	return group, nil
}

func (uc *SocialUsecase) GetGroup(id int) (*domain.Group, error) {
	group, err := uc.socialRepo.GetGroup(id)
	if err != nil {
		return nil, ErrGroupNotFound
	}
	return group, nil
}

func (uc *SocialUsecase) GetUserGroups(userID int) ([]domain.Group, error) {
	groups, err := uc.socialRepo.GetUserGroups(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}
	return groups, nil
}

func (uc *SocialUsecase) GetPublicGroups() ([]domain.Group, error) {
	groups, err := uc.socialRepo.GetPublicGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get public groups: %w", err)
	}
	return groups, nil
}

func (uc *SocialUsecase) JoinGroup(groupID, userID int) error {
	if err := uc.socialRepo.JoinGroup(groupID, userID, "member"); err != nil {
		return fmt.Errorf("failed to join group: %w", err)
	}
	return nil
}

func (uc *SocialUsecase) LeaveGroup(groupID, userID int) error {
	if err := uc.socialRepo.LeaveGroup(groupID, userID); err != nil {
		return fmt.Errorf("failed to leave group: %w", err)
	}
	return nil
}

func (uc *SocialUsecase) GetGroupMembers(groupID int) ([]domain.GroupMember, error) {
	members, err := uc.socialRepo.GetGroupMembers(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}
	return members, nil
}