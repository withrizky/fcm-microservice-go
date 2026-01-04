package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"fcm_microservice/internal/fcm"
	"fcm_microservice/internal/model"
	"fcm_microservice/internal/worker"
)

func main() {
	godotenv.Load()

	// 1. Setup FCM Client
	credFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credFile == "" {
		// Fallback ke nama file default jika di .env kosong
		credFile = "signal-app-17739-firebase-adminsdk-fbsvc-09ebee41d8.json"
	}

	fcmClient, err := fcm.NewClient(credFile)
	if err != nil {
		log.Fatalf("Gagal inisialisasi Firebase: %v", err)
	}

	// 2. Setup Worker Pool (Pure Go)
	// 50 Worker, Antrean Buffer 5000
	dispatcher := worker.NewDispatcher(50, 5000, fcmClient)
	dispatcher.Run()

	// 3. Setup Gin
	r := gin.New()
	r.Use(gin.Recovery())

	// Endpoint Send FCM
	r.POST("/send-fcm", func(c *gin.Context) {
		var req model.FcmPayload
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Payload tidak valid: " + err.Error()})
			return
		}

		// Masukkan ke antrean (Non-blocking)
		select {
		case dispatcher.JobQueue <- req:
			c.JSON(202, gin.H{"status": "queued", "message": "Notifikasi sedang diproses"})
		default:
			c.JSON(503, gin.H{"error": "Antrean penuh, server sibuk"})
		}
	})

	// 4. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("🔥 FCM Microservice running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	// 5. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Menyelesaikan antrean notifikasi...")
	dispatcher.Stop()
	log.Println("Server exiting")
}
