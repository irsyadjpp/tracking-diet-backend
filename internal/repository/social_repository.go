package repository

import (
	"github.com/irsyadjpp/tracking-diet-backend/internal/domain"
	"gorm.io/gorm"
)

type socialRepository struct {
	db *gorm.DB
}

func NewSocialRepository(db *gorm.DB) domain.SocialRepository {
	return &socialRepository{db: db}
}

// Connections
func (r *socialRepository) CreateConnection(connection *domain.Connection) error {
	return r.db.Create(connection).Error
}

func (r *socialRepository) GetConnection(userID, connectedID int) (*domain.Connection, error) {
	var connection domain.Connection
	if err := r.db.Where("(user_id = ? AND connected_id = ?) OR (user_id = ? AND connected_id = ?)", userID, connectedID, connectedID, userID).First(&connection).Error; err != nil {
		return nil, err
	}
	return &connection, nil
}

func (r *socialRepository) GetUserConnections(userID int) ([]domain.Connection, error) {
	var connections []domain.Connection
	if err := r.db.Where("(user_id = ? OR connected_id = ?) AND status = ?", userID, userID, domain.ConnectionStatusAccepted).Find(&connections).Error; err != nil {
		return nil, err
	}
	return connections, nil
}

func (r *socialRepository) GetPendingConnections(userID int) ([]domain.Connection, error) {
	var connections []domain.Connection
	if err := r.db.Where("connected_id = ? AND status = ?", userID, domain.ConnectionStatusPending).Find(&connections).Error; err != nil {
		return nil, err
	}
	return connections, nil
}

func (r *socialRepository) UpdateConnectionStatus(id int, status domain.ConnectionStatus) error {
	return r.db.Model(&domain.Connection{}).Where("id = ?", id).Update("status", status).Error
}

func (r *socialRepository) DeleteConnection(id int) error {
	return r.db.Delete(&domain.Connection{}, id).Error
}

// Posts
func (r *socialRepository) CreatePost(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *socialRepository) GetPost(id int) (*domain.Post, error) {
	var post domain.Post
	if err := r.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *socialRepository) GetUserPosts(userID int) ([]domain.Post, error) {
	var posts []domain.Post
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *socialRepository) GetFeedPosts(userID int) ([]domain.Post, error) {
	var posts []domain.Post
	// Get posts from user's connections
	var connectionIDs []int
	r.db.Model(&domain.Connection{}).
		Where("(user_id = ? OR connected_id = ?) AND status = ?", userID, userID, domain.ConnectionStatusAccepted).
		Pluck("user_id", &connectionIDs)
	
	// Add user's own posts
	connectionIDs = append(connectionIDs, userID)

	if err := r.db.Where("user_id IN ? AND is_public = ?", connectionIDs, true).Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *socialRepository) UpdatePost(post *domain.Post) error {
	return r.db.Save(post).Error
}

func (r *socialRepository) DeletePost(id int) error {
	return r.db.Delete(&domain.Post{}, id).Error
}

// Comments
func (r *socialRepository) CreateComment(comment *domain.Comment) error {
	if err := r.db.Create(comment).Error; err != nil {
		return err
	}
	// Update comment count
	return r.db.Model(&domain.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comments_count", gorm.Expr("comments_count + 1")).Error
}

func (r *socialRepository) GetPostComments(postID int) ([]domain.Comment, error) {
	var comments []domain.Comment
	if err := r.db.Where("post_id = ?", postID).Order("created_at ASC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *socialRepository) UpdateComment(comment *domain.Comment) error {
	return r.db.Save(comment).Error
}

func (r *socialRepository) DeleteComment(id int) error {
	var comment domain.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		return err
	}
	
	if err := r.db.Delete(&domain.Comment{}, id).Error; err != nil {
		return err
	}
	
	// Update comment count
	return r.db.Model(&domain.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comments_count", gorm.Expr("comments_count - 1")).Error
}

// Likes
func (r *socialRepository) CreateLike(postID, userID int) error {
	like := &domain.Like{
		PostID: postID,
		UserID: userID,
	}
	
	if err := r.db.Create(like).Error; err != nil {
		return err
	}
	
	// Update like count
	return r.db.Model(&domain.Post{}).Where("id = ?", postID).UpdateColumn("likes_count", gorm.Expr("likes_count + 1")).Error
}

func (r *socialRepository) DeleteLike(postID, userID int) error {
	if err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&domain.Like{}).Error; err != nil {
		return err
	}
	
	// Update like count
	return r.db.Model(&domain.Post{}).Where("id = ?", postID).UpdateColumn("likes_count", gorm.Expr("likes_count - 1")).Error
}

func (r *socialRepository) GetPostLikes(postID int) ([]domain.Like, error) {
	var likes []domain.Like
	if err := r.db.Where("post_id = ?", postID).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *socialRepository) CheckUserLiked(postID, userID int) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.Like{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Groups
func (r *socialRepository) CreateGroup(group *domain.Group) error {
	return r.db.Create(group).Error
}

func (r *socialRepository) GetGroup(id int) (*domain.Group, error) {
	var group domain.Group
	if err := r.db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *socialRepository) GetUserGroups(userID int) ([]domain.Group, error) {
	var groups []domain.Group
	if err := r.db.Joins("JOIN group_members ON groups.id = group_members.group_id").
		Where("group_members.user_id = ?", userID).
		Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *socialRepository) GetPublicGroups() ([]domain.Group, error) {
	var groups []domain.Group
	if err := r.db.Where("is_public = ?", true).Order("member_count DESC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *socialRepository) JoinGroup(groupID, userID int, role string) error {
	member := &domain.GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Role:    role,
	}
	
	if err := r.db.Create(member).Error; err != nil {
		return err
	}
	
	// Update member count
	return r.db.Model(&domain.Group{}).Where("id = ?", groupID).UpdateColumn("member_count", gorm.Expr("member_count + 1")).Error
}

func (r *socialRepository) LeaveGroup(groupID, userID int) error {
	if err := r.db.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&domain.GroupMember{}).Error; err != nil {
		return err
	}
	
	// Update member count
	return r.db.Model(&domain.Group{}).Where("id = ?", groupID).UpdateColumn("member_count", gorm.Expr("member_count - 1")).Error
}

func (r *socialRepository) GetGroupMembers(groupID int) ([]domain.GroupMember, error) {
	var members []domain.GroupMember
	if err := r.db.Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}