package internal

import (
	"fmt"
	"time"
)

// SaveMessage 保存聊天消息 - 核心函数
func SaveMessage(sessionID, role, content string) error {
	if sessionID == "" {
		sessionID = "default"
	}

	message := ChatMessage{
		SessionID: sessionID,
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}

	result := DB.Create(&message)
	if result.Error != nil {
		return fmt.Errorf("failed to save message: %v", result.Error)
	}

	return nil
}

// GetSessionMessages 获取会话的所有消息
func GetSessionMessages(sessionID string) ([]ChatMessage, error) {
	var messages []ChatMessage
	result := DB.Where("session_id = ?", sessionID).
		Order("timestamp ASC").
		Find(&messages)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get messages: %v", result.Error)
	}

	return messages, nil
}

// GetMessageCount 获取会话消息数量
func GetMessageCount(sessionID string) (int64, error) {
	var count int64
	result := DB.Model(&ChatMessage{}).Where("session_id = ?", sessionID).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to count messages: %v", result.Error)
	}
	return count, nil
}

// GetUserMessages 获取用户消息（用于生成画像）
func GetUserMessages(sessionID string) ([]string, error) {
	var messages []ChatMessage
	result := DB.Where("session_id = ? AND role = ?", sessionID, "user").
		Order("timestamp ASC").
		Select("content").
		Find(&messages)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user messages: %v", result.Error)
	}

	userMessages := make([]string, len(messages))
	for i, msg := range messages {
		userMessages[i] = msg.Content
	}

	return userMessages, nil
}

// GetAllSessionIDs 获取所有会话ID
func GetAllSessionIDs() ([]string, error) {
	var sessionIDs []string
	result := DB.Model(&ChatMessage{}).
		Distinct("session_id").
		Pluck("session_id", &sessionIDs)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get session IDs: %v", result.Error)
	}

	return sessionIDs, nil
}

// SaveProfile 保存用户画像
func SaveProfile(sessionID, profile string) error {
	profileRecord := UserProfile{
		SessionID: sessionID,
		Profile:   profile,
	}

	// 使用GORM的Upsert功能 - 存在则更新，不存在则创建
	result := DB.Where("session_id = ?", sessionID).
		Assign(UserProfile{Profile: profile}).
		FirstOrCreate(&profileRecord)

	if result.Error != nil {
		return fmt.Errorf("failed to save profile: %v", result.Error)
	}

	return nil
}

// GetProfile 获取用户画像
func GetProfile(sessionID string) (*UserProfile, error) {
	var profile UserProfile
	result := DB.Where("session_id = ?", sessionID).First(&profile)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get profile: %v", result.Error)
	}
	return &profile, nil
}

// SearchMessages 搜索消息内容
func SearchMessages(sessionID, keyword string) ([]ChatMessage, error) {
	var messages []ChatMessage
	result := DB.Where("session_id = ? AND content LIKE ?", sessionID, "%"+keyword+"%").
		Order("timestamp DESC").
		Find(&messages)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to search messages: %v", result.Error)
	}

	return messages, nil
}

// GetRecentMessages 获取最近的消息
func GetRecentMessages(sessionID string, limit int) ([]ChatMessage, error) {
	var messages []ChatMessage
	result := DB.Where("session_id = ?", sessionID).
		Order("timestamp DESC").
		Limit(limit).
		Find(&messages)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recent messages: %v", result.Error)
	}

	return messages, nil
}

// GetMessagesByDateRange 根据时间范围获取消息
func GetMessagesByDateRange(sessionID string, startTime, endTime time.Time) ([]ChatMessage, error) {
	var messages []ChatMessage
	result := DB.Where("session_id = ? AND timestamp BETWEEN ? AND ?", sessionID, startTime, endTime).
		Order("timestamp ASC").
		Find(&messages)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get messages by date range: %v", result.Error)
	}

	return messages, nil
}

// DeleteSession 删除整个会话的所有消息
func DeleteSession(sessionID string) error {
	// 删除消息
	result := DB.Where("session_id = ?", sessionID).Delete(&ChatMessage{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete messages: %v", result.Error)
	}

	// 删除画像
	DB.Where("session_id = ?", sessionID).Delete(&UserProfile{})

	return nil
}

// GetSessionStats 获取会话统计信息
func GetSessionStats(sessionID string) (map[string]interface{}, error) {
	var userCount, assistantCount int64
	var firstMessage, lastMessage time.Time

	// 统计用户消息数
	DB.Model(&ChatMessage{}).Where("session_id = ? AND role = ?", sessionID, "user").Count(&userCount)

	// 统计助手消息数
	DB.Model(&ChatMessage{}).Where("session_id = ? AND role = ?", sessionID, "assistant").Count(&assistantCount)

	// 获取第一条和最后一条消息时间
	DB.Model(&ChatMessage{}).Where("session_id = ?", sessionID).
		Order("timestamp ASC").Limit(1).Pluck("timestamp", &firstMessage)

	DB.Model(&ChatMessage{}).Where("session_id = ?", sessionID).
		Order("timestamp DESC").Limit(1).Pluck("timestamp", &lastMessage)

	stats := map[string]interface{}{
		"user_message_count":      userCount,
		"assistant_message_count": assistantCount,
		"total_message_count":     userCount + assistantCount,
		"first_message_time":      firstMessage,
		"last_message_time":       lastMessage,
	}

	return stats, nil
}
