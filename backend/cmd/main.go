package main

import (
	"log"
	"os"

	"github.com/Adityoexs/contract-employee/backend/internal/db"
	"github.com/Adityoexs/contract-employee/backend/internal/handler"
	"github.com/Adityoexs/contract-employee/backend/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Database
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()
	log.Println("Database connected successfully")

	// Repository & Handler
	repo := repository.NewKaryawanRepository(dbConn)
	h := handler.NewKaryawanHandler(repo)

	// Router
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// Routes
	api := r.Group("/api")
	{
		karyawan := api.Group("/karyawan")
		{
			karyawan.GET("", h.GetAll)
			karyawan.GET("/:id", h.GetByID)
			karyawan.POST("", h.Create)
			karyawan.PUT("/:id", h.Update)
			karyawan.DELETE("/:id", h.Delete)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
