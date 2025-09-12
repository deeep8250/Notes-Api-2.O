package main

import (
	"fmt"
	"log"
	"pr01/config"
	"pr01/db"
	"pr01/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	cfg := config.Load()
	db.DbInit(*cfg)
	r := gin.Default()
	routes.Routes(r)

	addr := fmt.Sprintf(":%d", cfg.Port)

	if err := r.Run(addr); err != nil {
		log.Fatal("server failed:", err)
	}

}

//ererrrr
