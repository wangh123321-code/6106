package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tt-tournament/backend/internal/middleware"
	"github.com/tt-tournament/backend/internal/model"
	"github.com/tt-tournament/backend/internal/service"
	"github.com/tt-tournament/backend/pkg/response"
)

type MatchHandler struct {
	svc *service.MatchService
}

func NewMatchHandler(svc *service.MatchService) *MatchHandler {
	return &MatchHandler{svc: svc}
}

func (h *MatchHandler) RecordScore(c *gin.Context) {
	matchID := c.Param("match_id")
	refereeID := c.GetString("user_id")

	var req struct {
		Games []model.GameScore `json:"games" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	input := service.ScoreInput{
		MatchID: matchID,
		Games:   req.Games,
	}
	match, err := h.svc.RecordScore(c.Request.Context(), refereeID, input)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, match)
}

func (h *MatchHandler) Walkover(c *gin.Context) {
	matchID := c.Param("match_id")
	var req struct {
		PlayerID string `json:"player_id" binding:"required"`
		Reason   string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	match, err := h.svc.Walkover(c.Request.Context(), matchID, req.PlayerID, req.Reason)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, match)
}

func (h *MatchHandler) OverrideScore(c *gin.Context) {
	matchID := c.Param("match_id")
	operatorID := c.GetString("user_id")
	operatorName := c.GetString("username")

	var req struct {
		Games []model.GameScore `json:"games" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	input := service.ScoreInput{
		MatchID: matchID,
		Games:   req.Games,
	}
	match, err := h.svc.OverrideScore(c.Request.Context(), operatorID, operatorName, matchID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, match)
}

func (h *MatchHandler) GetByID(c *gin.Context) {
	matchID := c.Param("match_id")
	match, err := h.svc.GetByID(c.Request.Context(), matchID)
	if err != nil {
		response.NotFound(c, "比赛不存在")
		return
	}
	response.Success(c, match)
}

func (h *MatchHandler) AssignReferee(c *gin.Context) {
	matchID := c.Param("match_id")
	var req struct {
		RefereeID string `json:"referee_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	err := h.svc.AssignReferee(c.Request.Context(), matchID, req.RefereeID)
	if err != nil {
		response.InternalError(c, "指派裁判失败")
		return
	}
	response.Success(c, nil)
}

func (h *MatchHandler) ListByReferee(c *gin.Context) {
	refereeID := c.GetString("user_id")
	list, err := h.svc.ListByReferee(c.Request.Context(), refereeID)
	if err != nil {
		response.InternalError(c, "获取比赛列表失败")
		return
	}
	response.Success(c, list)
}

func (h *MatchHandler) GetAuditLogs(c *gin.Context) {
	targetType := c.Query("target_type")
	targetID := c.Query("target_id")
	logs, err := h.svc.GetAuditLogs(c.Request.Context(), targetType, targetID)
	if err != nil {
		response.InternalError(c, "获取审计日志失败")
		return
	}
	response.Success(c, logs)
}

func (h *MatchHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("", middleware.AuthRequired())
	{
		auth.POST("/matches/:match_id/score", middleware.RoleRequired("referee", "committee"), h.RecordScore)
		auth.POST("/matches/:match_id/walkover", middleware.RoleRequired("committee", "referee"), h.Walkover)
		auth.PUT("/matches/:match_id/override", middleware.RoleRequired("committee"), h.OverrideScore)
		auth.GET("/matches/:match_id", h.GetByID)
		auth.POST("/matches/:match_id/assign-referee", middleware.RoleRequired("committee"), h.AssignReferee)
		auth.GET("/referee/matches", middleware.RoleRequired("referee"), h.ListByReferee)
		auth.GET("/audit-logs", middleware.RoleRequired("committee"), h.GetAuditLogs)
	}
}
