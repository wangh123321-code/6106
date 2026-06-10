package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tt-tournament/backend/internal/service"
	"github.com/tt-tournament/backend/pkg/response"
)

type RankingHandler struct {
	svc *service.RankingService
}

func NewRankingHandler(svc *service.RankingService) *RankingHandler {
	return &RankingHandler{svc: svc}
}

func (h *RankingHandler) GetMedalBoard(c *gin.Context) {
	response.Success(c, gin.H{"message": "奖牌榜数据需根据赛事状态动态计算"})
}

func (h *RankingHandler) ExportCertificate(c *gin.Context) {
	playerName := c.Query("player_name")
	eventName := c.Query("event_name")
	medal := c.Query("medal")
	tournamentName := c.Query("tournament_name")
	date := c.Query("date")

	pdf, err := service.GenerateCertificatePDF(playerName, eventName, medal, tournamentName, date)
	if err != nil {
		response.InternalError(c, "生成证书失败")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=certificate.pdf")
	c.Data(200, "application/pdf", pdf)
}

func (h *RankingHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/tournaments/:id/medal-board", h.GetMedalBoard)
	rg.GET("/certificate/export", h.ExportCertificate)
}
