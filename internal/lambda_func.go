package internal

import (
	"context"
	"fmt"
	"strings"
)

var sessions = make(map[string][]map[string]interface{})

// newLambda component initialization function of node 'InputHandler' in graph 'awesomeeino'
func newLambda(_ctx context.Context, input UserInput, opts ...any) (output ProcessedInput, err error) {
	//处理输入
	//主要就是检查id还有查询聊天记录条数 ==>生成用户画像用的
	if input.SessionID == "" {
		input.SessionID = "default"
	}
	messageCount, err := GetMessageCount(input.SessionID)
	return ProcessedInput{
		Message:      input.Message,
		MessageCount: int(messageCount),
		SessionID:    input.SessionID,
		Timestamp:    "",
	}, nil
}

// newLambda1 component initialization function of node 'extractMessages' in graph 'awesomeeino'
func newLambda1(_ctx context.Context, input UserInput) (output any, err error) {
	//这个是用来提取用户消息=>生成画像
	userMessages, err := GetUserMessages(input.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user messages: %v", err)
	}
	return map[string]string{
		//拼接
		"user_messages": strings.Join(userMessages, "\n"),
		"session_ID":    input.SessionID,
	}, nil
}

// newLambda2 component initialization function of node 'Lambda5' in graph 'awesomeeino'
func newLambda2(_ctx context.Context, input any) (output any, err error) {
	//处理画像
	data, ok := input.(map[string]interface{})
	if !ok {
		return "Profile generation completed", nil
	}
	//把转化为string 吗，感觉会不会被多一此举
	SessionID, _ := data["session_id"].(string)

	profile, _ := data["content"].(string)

	SaveProfile(SessionID, profile)
	return fmt.Sprintf("用户画像以及保存成功"), nil

}

// newLambda3 component initialization function of node 'OutputHandler' in graph 'awesomeeino'
func newLambda3(_ctx context.Context, input ChatResponse) (output ProcessedResponse, err error) {
	//保存对话历史
	SaveMessage(input.SessionID, "assistant", input.Content)
	newCount := input.MessageCount + 1
	trigger_profile := newCount%10 == 0
	return ProcessedResponse{
		Content:        input.Content,
		MessageCount:   newCount,
		SessionID:      input.SessionID,
		TriggerProfile: trigger_profile,
	}, nil
}

// 保存会话到文件
