package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	config "github.com/maheshrc27/postflow/configs"
	"github.com/maheshrc27/postflow/internal/api/handlers"
	"github.com/maheshrc27/postflow/internal/api/middleware"
	"github.com/maheshrc27/postflow/internal/repository"
	"github.com/maheshrc27/postflow/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Failed to load environment variables", err)
	}

	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.PostgresURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer closeDB(db)

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		BodyLimit:    100 * 1024 * 1024, // 100 MB
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://locahost:3000, http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	userRepo := repository.NewUserRepository(db)
	mediaAssetRepo := repository.NewMediaAssetRepository(db)
	creditsRepo := repository.NewCreditsRepository(db)

	authService := service.NewAuthService(*cfg, userRepo, creditsRepo)
	userService := service.NewUserService(userRepo)
	creditsService := service.NewCreditsService(creditsRepo)
	videoService := service.NewVideoService(creditsRepo, mediaAssetRepo, *cfg)
	paymentService := service.NewPaymentService(*cfg, userRepo, creditsRepo)

	auth := handlers.NewAuthHandler(*cfg, authService)
	app.Get("/login", auth.Login)
	app.Get("/login/callback", auth.LoginCallbackHandler)

	payment := handlers.NewPaymentHandler(paymentService)
	app.Post("/payment/webhook", payment.PaymentWebhook)

	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))

	user := handlers.NewUserHandler(userService, *cfg)
	api.Get("/user/info", user.GetUserInfo)
	api.Post("/user/delete", user.DeleteAccount)

	credits := handlers.NewCreditsHandler(creditsService)
	api.Get("/credits", credits.GetCredits)

	video := handlers.NewVideoHandler(videoService)
	api.Get("/videos", video.GetVideos)
	api.Post("/generate", video.CreateVideo)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Println("Server is running on http://localhost:3000")

	gracefulShutdown(app, db)
}

func closeDB(db *sql.DB) {
	fmt.Fprint(os.Stdout, "Closing database connection... ")
	if err := db.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close database: %v", err)
		return
	}
	fmt.Fprintln(os.Stdout, "Done")
}

func gracefulShutdown(app *fiber.App, db *sql.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Failed to shut down server: %v", err)
	}

	closeDB(db)
	log.Println("Server shutdown complete.")
}
