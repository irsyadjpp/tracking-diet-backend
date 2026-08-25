package domain

import "time"

// ConnectionType represents different types of social connections
type ConnectionType string

const (
	ConnectionTypeFollow   ConnectionType = "follow"
	ConnectionTypeFriend   ConnectionType = "friend"
	ConnectionTypeMentor   ConnectionType = "mentor"
	ConnectionTypeMentee   ConnectionType = "mentee"
)

// ConnectionStatus represents the status of a connection
type ConnectionStatus string

const (
	ConnectionStatusPending  ConnectionStatus = "pending"
	ConnectionStatusAccepted ConnectionStatus = "accepted"
	ConnectionStatusRejected ConnectionStatus = "rejected"
	ConnectionStatusBlocked  ConnectionStatus = "blocked"
)

// Connection represents a social connection between users
type Connection struct {
	ID          int             `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int             `gorm:"not null;index" json:"user_id"`
	ConnectedID int             `gorm:"not null;index" json:"connected_id"`
	Type        ConnectionType  `gorm:"size:20;not null" json:"type"`
	Status      ConnectionStatus `gorm:"size:20;default:pending" json:"status"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Connection model
func (Connection) TableName() string {
	return "connections"
}

// PostType represents different types of community posts
type PostType string

const (
	PostTypeAchievement   PostType = "achievement"
	PostTypeProgress      PostType = "progress"
	PostTypeQuestion      PostType = "question"
	PostTypeMotivation    PostType = "motivation"
	PostTypeRecipe        PostType = "recipe"
	PostTypeChallenge     PostType = "challenge"
)

// Post represents a community post
type Post struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"not null;index" json:"user_id"`
	Type        PostType  `gorm:"size:50;not null" json:"type"`
	Title       string    `gorm:"size:200" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	ImageURL    string    `gorm:"size:500" json:"image_url"`
	IsPublic    bool      `gorm:"default:true" json:"is_public"`
	LikesCount  int       `gorm:"default:0" json:"likes_count"`
	CommentsCount int     `gorm:"default:0" json:"comments_count"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Post model
func (Post) TableName() string {
	return "posts"
}

// Comment represents a comment on a post
type Comment struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    int       `gorm:"not null;index" json:"post_id"`
	UserID    int       `gorm:"not null;index" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Comment model
func (Comment) TableName() string {
	return "comments"
}

// Like represents a like on a post
type Like struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    int       `gorm:"not null;index" json:"post_id"`
	UserID    int       `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for Like model
func (Like) TableName() string {
	return "likes"
}

// Group represents a community group
type Group struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	ImageURL    string    `gorm:"size:500" json:"image_url"`
	IsPublic    bool      `gorm:"default:true" json:"is_public"`
	MemberCount int       `gorm:"default:0" json:"member_count"`
	CreatedBy   int       `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Group model
func (Group) TableName() string {
	return "groups"
}

// GroupMember represents a user's membership in a group
type GroupMember struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID    int       `gorm:"not null;index" json:"group_id"`
	UserID     int       `gorm:"not null;index" json:"user_id"`
	Role       string    `gorm:"size:20;default:member" json:"role"` // admin, moderator, member
	JoinedAt   time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

// TableName specifies the table name for GroupMember model
func (GroupMember) TableName() string {
	return "group_members"
}

// SocialRepository defines operations for managing social features
type SocialRepository interface {
	// Connections
	CreateConnection(connection *Connection) error
	GetConnection(userID, connectedID int) (*Connection, error)
	GetUserConnections(userID int) ([]Connection, error)
	GetPendingConnections(userID int) ([]Connection, error)
	UpdateConnectionStatus(id int, status ConnectionStatus) error
	DeleteConnection(id int) error

	// Posts
	CreatePost(post *Post) error
	GetPost(id int) (*Post, error)
	GetUserPosts(userID int) ([]Post, error)
	GetFeedPosts(userID int) ([]Post, error)
	UpdatePost(post *Post) error
	DeletePost(id int) error

	// Comments
	CreateComment(comment *Comment) error
	GetPostComments(postID int) ([]Comment, error)
	UpdateComment(comment *Comment) error
	DeleteComment(id int) error

	// Likes
	CreateLike(postID, userID int) error
	DeleteLike(postID, userID int) error
	GetPostLikes(postID int) ([]Like, error)
	CheckUserLiked(postID, userID int) (bool, error)

	// Groups
	CreateGroup(group *Group) error
	GetGroup(id int) (*Group, error)
	GetUserGroups(userID int) ([]Group, error)
	GetPublicGroups() ([]Group, error)
	JoinGroup(groupID, userID int, role string) error
	LeaveGroup(groupID, userID int) error
	GetGroupMembers(groupID int) ([]GroupMember, error)
}