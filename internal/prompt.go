package internal

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type ChatTemplateConfig struct {
	FormatType schema.FormatType
	Templates  []schema.MessagesTemplate
}

// newChatTemplate component initialization function of node 'ChatTemplate2' in graph 'awesomeeino'
func newChatTemplate(_ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	//用于生成画像
	// TODO Modify component configuration here.
	ctp = prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是用户行为分析专家，请根据用户的聊天历史生成用户画像。"),
		schema.UserMessage(`请分析以下用户消息，生成用户画像：

{user_messages}

请从以下维度分析：
1. 年龄段和性别推测
2. 兴趣爱好特征  
3. 性格特点
4. 沟通风格

请用自然语言描述。`),
	)
	return ctp, nil
}

// newChatTemplate1 component initialization function of node 'ChatTemplate1' in graph 'awesomeeino'
func newChatTemplate1(_ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	//聊天
	// TODO Modify component configuration here.
	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个有用的AI助手，请自然地与用户对话。"),
		schema.UserMessage("{message}"),
	)
	return template, nil
}
