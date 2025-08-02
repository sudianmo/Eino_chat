package main

import (
	"AgentWithGraph/internal"
	"log"
)

func main() {
	log.Println("Starting AI Chat Service...")
	// 初始化数据库
	if err := internal.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer internal.CloseDatabase()
	internal.StartServer()
}
