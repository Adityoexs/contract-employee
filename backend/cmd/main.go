package main

import (
	"log"
	"os"

	"github.com/Adityoexs/contract-employee/backend/internal/db"
	"github.com/Adityoexs/contract-employee/backend/internal/handler"
	"github.com/Adityoexs/contract-employee/backend/internal/repository"
	"github.com/Adityoexs/contract-employee/backend/internal/service"
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

	var currentDB, currentSchema, searchPath string
	if err := dbConn.QueryRow("SELECT current_database(), current_schema(), current_setting('search_path')").Scan(&currentDB, &currentSchema, &searchPath); err != nil {
		log.Fatalf("Failed to inspect DB: %v", err)
	}
	log.Printf("MAIN connected to database=%s schema=%s search_path=%s", currentDB, currentSchema, searchPath)

	// Repository, Service & Handler
	repo := repository.NewKaryawanRepository(dbConn)
	svc := service.NewKaryawanService(repo)
	h := handler.NewKaryawanHandler(svc)

	// Router
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
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
			karyawan.POST("/import", h.ImportExcel)
			karyawan.GET("/import/template", h.DownloadTemplate)
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
	log.Println("STARTUP: main entered")
	log.Println("DB_NAME:", os.Getenv("DB_NAME"))
	log.Println("DB_HOST:", os.Getenv("DB_HOST"))
	log.Println("DB_PORT:", os.Getenv("DB_PORT"))
	log.Println("DB_USER:", os.Getenv("DB_USER"))
}
