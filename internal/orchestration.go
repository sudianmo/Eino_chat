package internal

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func Buildawesomeeino(ctx context.Context) (r compose.Runnable[string, string], err error) {
	const (
		ChatModel2      = "ChatModel2"
		ChatTemplate2   = "ChatTemplate2"
		InputHandler    = "InputHandler"
		extractMessages = "extractMessages"
		Lambda5         = "Lambda5"
		ChatModel1      = "ChatModel1"
		OutputHandler   = "OutputHandler"
		ChatTemplate1   = "ChatTemplate1"
	)
	g := compose.NewGraph[string, string]()
	chatModel2KeyOfChatModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(ChatModel2, chatModel2KeyOfChatModel)
	chatTemplate2KeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate2, chatTemplate2KeyOfChatTemplate)
	_ = g.AddLambdaNode(InputHandler, compose.InvokableLambdaWithOption(newLambda))
	_ = g.AddLambdaNode(extractMessages, compose.InvokableLambda(newLambda1))
	_ = g.AddLambdaNode(Lambda5, compose.InvokableLambda(newLambda2))
	chatModel1KeyOfChatModel, err := newChatModel1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(ChatModel1, chatModel1KeyOfChatModel)
	_ = g.AddLambdaNode(OutputHandler, compose.InvokableLambda(newLambda3))
	chatTemplate1KeyOfChatTemplate, err := newChatTemplate1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate1, chatTemplate1KeyOfChatTemplate)
	_ = g.AddEdge(Lambda5, compose.END)
	_ = g.AddEdge(OutputHandler, compose.END)
	_ = g.AddEdge(ChatTemplate2, ChatModel2)
	_ = g.AddEdge(ChatModel2, Lambda5)
	_ = g.AddEdge(extractMessages, ChatTemplate2)
	_ = g.AddEdge(InputHandler, ChatTemplate1)
	_ = g.AddEdge(ChatTemplate1, ChatModel1)
	_ = g.AddEdge(ChatModel1, OutputHandler)
	_ = g.AddBranch(compose.START, compose.NewGraphBranch(newBranch, map[string]bool{InputHandler: true, extractMessages: true}))
	r, err = g.Compile(ctx, compose.WithGraphName("awesomeeino"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
