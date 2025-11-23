package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Traunin/review-assigner/internal/api/handlers"
	"github.com/Traunin/review-assigner/internal/application/services"
	"github.com/Traunin/review-assigner/internal/config"
	domainservices "github.com/Traunin/review-assigner/internal/domain/services"
	"github.com/Traunin/review-assigner/internal/infrastructure/db/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	config := config.Load()

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser(),
		config.DBPassword(),
		config.DBHost(),
		config.DBPort(),
		config.DBName(),
	)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	db := postgres.NewDB(pool)
	log.Println("Successfully connected to database")

	userRepo := postgres.NewUserRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	prRepo := postgres.NewPullRequestRepository(db)

	assignmentService := domainservices.NewReviewerAssignmentService(
		userRepo,
		prRepo,
		teamRepo,
	)

	teamService := services.NewTeamService(teamRepo, userRepo)
	prService := services.NewPullRequestService(prRepo, assignmentService)

	server := handlers.NewServer(
		teamService,
		prService,
		userRepo,
		teamRepo,
		prRepo,
	)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	registerRoutes(e, server)

	port := config.Port()
	log.Printf("Starting server on :%s", port)
	if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(e *echo.Echo, server *handlers.Server) {
	e.POST("/team/add", server.PostTeamAdd)
	e.GET("/team/get", server.GetTeamGet)

	e.POST("/users/setIsActive", server.PostUsersSetIsActive)
	e.GET("/users/getReview", server.GetUsersGetReview)

	e.POST("/pullRequest/create", server.PostPullRequestCreate)
	e.POST("/pullRequest/merge", server.PostPullRequestMerge)
	e.POST("/pullRequest/reassign", server.PostPullRequestReassign)

	// healthcheck
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
}
