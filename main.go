package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"docker-workshop-assesment-grader/internal/database"
	"docker-workshop-assesment-grader/internal/handlers"
	"docker-workshop-assesment-grader/internal/sse"
)

//go:embed all:frontend/dist
var frontendDist embed.FS

func main() {
	db, err := database.Init("students.db")
	if err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	hub := sse.NewHub()

	studentHandler := &handlers.StudentHandler{DB: db}
	registerHandler := &handlers.RegisterHandler{DB: db, Hub: hub}
	approvalHandler := &handlers.ApprovalHandler{DB: db}
	notifyHandler := &handlers.NotifyHandler{DB: db, Hub: hub}
	eventsHandler := &handlers.EventsHandler{Hub: hub}

	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/students", studentHandler.CreateStudent)
		api.GET("/students", studentHandler.ListStudents)
		api.GET("/students/:id", studentHandler.GetStudent)
		api.PUT("/students/:id", studentHandler.UpdateStudent)
		api.DELETE("/students/:id", studentHandler.DeleteStudent)

		api.POST("/register", registerHandler.Register)
		api.POST("/notify", notifyHandler.Notify)
		api.GET("/events", eventsHandler.Stream)

		api.POST("/registrations/:id/approve", approvalHandler.ApproveOne)
		api.POST("/registrations/approve-all", approvalHandler.ApproveAll)
	}

	registerFrontendRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}

func registerFrontendRoutes(router *gin.Engine) {
	dist, err := fs.Sub(frontendDist, "frontend/dist")
	if err != nil {
		log.Printf("frontend dist not available in embed fs: %v", err)
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "frontend build not found"})
		})
		return
	}

	fileServer := http.FileServer(http.FS(dist))
	router.NoRoute(func(c *gin.Context) {
		requestPath := filepath.Clean(c.Request.URL.Path)
		if requestPath != "." && requestPath != "/" {
			if file, fileErr := dist.Open(requestPath[1:]); fileErr == nil {
				file.Close()
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		index, openErr := dist.Open("index.html")
		if openErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "index.html not found"})
			return
		}
		defer index.Close()

		stat, statErr := index.Stat()
		if statErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read index metadata"})
			return
		}

		c.DataFromReader(http.StatusOK, stat.Size(), "text/html; charset=utf-8", index, nil)
	})
}
