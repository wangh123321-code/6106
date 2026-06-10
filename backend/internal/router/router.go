package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tt-tournament/backend/internal/handler"
)

func Setup(
	userH *handler.UserHandler,
	tournamentH *handler.TournamentHandler,
	eventH *handler.EventHandler,
	registrationH *handler.RegistrationHandler,
	bracketH *handler.BracketHandler,
	matchH *handler.MatchHandler,
	rankingH *handler.RankingHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	{
		userH.RegisterRoutes(api)
		tournamentH.RegisterRoutes(api)
		eventH.RegisterRoutes(api)
		registrationH.RegisterRoutes(api)
		bracketH.RegisterRoutes(api)
		matchH.RegisterRoutes(api)
		rankingH.RegisterRoutes(api)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
