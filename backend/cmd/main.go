package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/tt-tournament/backend/internal/handler"
	"github.com/tt-tournament/backend/internal/middleware"
	"github.com/tt-tournament/backend/internal/router"
	"github.com/tt-tournament/backend/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://ttadmin:ttsecret2026@localhost:27017/tournament?authSource=admin"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "changeme-in-production-32chars!!"
	}
	middleware.SetJWTSecret(jwtSecret)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("MongoDB 连接失败: %v", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB Ping 失败: %v", err)
	}
	log.Println("MongoDB 连接成功")

	db := client.Database("tournament")

	userSvc := service.NewUserService(db)
	tournamentSvc := service.NewTournamentService(db)
	eventSvc := service.NewEventService(db)
	registrationSvc := service.NewRegistrationService(db, userSvc)
	bracketSvc := service.NewBracketService(db, registrationSvc, eventSvc, userSvc)
	matchSvc := service.NewMatchService(db, userSvc, bracketSvc)
	rankingSvc := service.NewRankingService(eventSvc, userSvc)
	appealSvc := service.NewAppealService(db, matchSvc)

	userH := handler.NewUserHandler(userSvc)
	tournamentH := handler.NewTournamentHandler(tournamentSvc, appealSvc)
	eventH := handler.NewEventHandler(eventSvc)
	registrationH := handler.NewRegistrationHandler(registrationSvc)
	bracketH := handler.NewBracketHandler(bracketSvc)
	matchH := handler.NewMatchHandler(matchSvc)
	rankingH := handler.NewRankingHandler(rankingSvc)
	appealH := handler.NewAppealHandler(appealSvc)

	r := router.Setup(userH, tournamentH, eventH, registrationH, bracketH, matchH, rankingH, appealH)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("服务启动于 :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
