package internal

import (
	"context"
	"fmt"
	"time"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type ChatRequest struct {
	Message   string `json:"message"`
	SessionID string `json:"session_id"`
}
type ChatApiResponse struct {
	Content      string `json:"content"`
	MessageCount int    `json:"message_count"`
	SessionID    string `json:"session_id"`
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
}

func StartServer() {
	godotenv.Load()
	r := gin.Default()

	api := r.Group("api")
	{
		api.POST("/chat", handleChat)
		api.GET("/messages", getMessages)
		api.GET("/profile/:id", getProfile)
		api.GET("/sessions", listSessions)
	}
	//r.Static("/static", "./web")
	//r.LoadHTMLGlob("web/*.html")
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	log.Println("Server starting on :8080")
	r.Run(":8080")
}

func handleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, ApiResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// 转换为UserInput给图使用
	userInput := UserInput{
		SessionID: req.SessionID,
		Message:   req.Message,
	}
	SaveMessage(req.SessionID, "user", req.Message)
	ctx := context.Background()
	runnable, err := Buildawesomeeino(ctx)
	result, err := runnable.Invoke(ctx, userInput.Message) // 传入string
	if err != nil {
		log.Printf("Graph invoke error: %v", err)
		c.JSON(http.StatusInternalServerError, ChatApiResponse{
			Success: false,
			Error:   "Graph processing failed: " + err.Error(),
		})
		return
	}

	c.JSON(200, ChatApiResponse{
		Content:      result,
		SessionID:    req.SessionID,
		MessageCount: len(sessions[req.SessionID]),
		Success:      true,
	})
}

func getMessages(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		sessionID = "default"
	}

	messages, err := GetSessionMessages(sessionID)
	if err != nil {
		c.JSON(200, gin.H{
			"success": true,
			"data":    []interface{}{}, // 返回空数组而不是错误
		})
		return
	}
	var result []map[string]interface{}
	for _, msg := range messages {
		sender := "assistant"
		if msg.Role == "user" {
			sender = "user"
		}
		result = append(result, map[string]interface{}{
			"id":        fmt.Sprintf("%d", msg.ID),
			"content":   msg.Content,
			"sender":    sender,
			"timestamp": msg.Timestamp.UnixMilli(),
			"status":    "sent",
		})
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    result,
	})
}

//然后从文件获取
// 从文件获取

func getProfile(c *gin.Context) {
	sessionID := c.Param("id")
	filename := filepath.Join("data/profiles", sessionID+"_profile.json")

	data, err := os.ReadFile(filename)
	if err == nil {
		var profile map[string]interface{}
		json.Unmarshal(data, &profile)
		//为什么是unmarshal
		c.JSON(http.StatusOK, gin.H{
			"profile": profile,
			"success": true,
		})
		return
	}
	c.JSON(404, gin.H{
		"success": false,
		"error":   "Profile not found",
	})
}

func getSession(c *gin.Context) {
	sessionID := c.Param("id")
	//先从内存获取
	if data, exists := sessions[sessionID]; exists {
		content, _ := json.Marshal(data)

		c.JSON(200, ChatApiResponse{
			Content:      string(content),
			SessionID:    sessionID,
			MessageCount: len(data),
			Success:      true,
		})
		return
	}
	//然后从文件获取
	// 从文件获取
	filename := filepath.Join("data/sessions", sessionID+".json")
	if data, err := os.ReadFile(filename); err == nil {
		var messages []map[string]interface{}
		err := json.Unmarshal(data, &messages)
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   "Failed to unmarshal session data",
			})
			return
		}
		sessions[sessionID] = messages

		c.JSON(200, gin.H{
			"session_id": sessionID,
			"messages":   messages,
			"success":    true,
		})
		return
	}

	c.JSON(404, gin.H{
		"success": false,
		"error":   "Session not found",
	})
}

// 列出所有会话
func listSessions(c *gin.Context) {
	sessionList := make([]string, 0, len(sessions))
	for sessionID := range sessions {
		sessionList = append(sessionList, sessionID)
	}

	c.JSON(200, gin.H{
		"success":  true,
		"sessions": sessionList,
	})
}
