package internal

import "time"

type UserInput struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

type ProcessedInput struct {
	SessionID    string `json:"session_id"`
	Message      string `json:"message"`
	Timestamp    string `json:"timestamp"`
	MessageCount int    `json:"message_count"`
}

type ChatResponse struct {
	Content      string `json:"content"`
	SessionID    string `json:"session_id"`
	MessageCount int    `json:"message_count"`
}

type ProcessedResponse struct {
	Content        string `json:"content"`
	SessionID      string `json:"session_id"`
	MessageCount   int    `json:"message_count"`
	TriggerProfile bool   `json:"trigger_profile"`
}
type ChatMessage struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID string    `gorm:"index;size:255;not null" json:"session_id"`
	Role      string    `gorm:"size:20;not null" json:"role"` // user, assistant
	Content   string    `gorm:"type:longtext;not null" json:"content"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// UserProfile 用户画像表 - 独立简单表
type UserProfile struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID string    `gorm:"uniqueIndex;size:255;not null" json:"session_id"`
	Profile   string    `gorm:"type:longtext;not null" json:"profile"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// 表名配置
func (ChatMessage) TableName() string {
	return "chat_messages"
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
