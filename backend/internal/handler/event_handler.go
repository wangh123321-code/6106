package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tt-tournament/backend/internal/middleware"
	"github.com/tt-tournament/backend/internal/service"
	"github.com/tt-tournament/backend/pkg/response"
)

type EventHandler struct {
	svc *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) Create(c *gin.Context) {
	tournamentID := c.Param("tournament_id")
	var req struct {
		Name       string `json:"name" binding:"required"`
		Type       string `json:"type" binding:"required,oneof=ms ws xd"`
		DrawMethod string `json:"draw_method" binding:"required,oneof=snake random"`
		SeedCount  int    `json:"seed_count"`
		BracketSize int   `json:"bracket_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	e, err := h.svc.Create(c.Request.Context(), tournamentID, req.Name, req.Type, req.DrawMethod, req.SeedCount, req.BracketSize)
	if err != nil {
		response.InternalError(c, "创建项目失败")
		return
	}
	response.Created(c, e)
}

func (h *EventHandler) List(c *gin.Context) {
	tournamentID := c.Param("tournament_id")
	list, err := h.svc.ListByTournament(c.Request.Context(), tournamentID)
	if err != nil {
		response.InternalError(c, "获取项目列表失败")
		return
	}
	response.Success(c, list)
}

func (h *EventHandler) CloseRegistration(c *gin.Context) {
	id := c.Param("event_id")
	err := h.svc.CloseRegistration(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "关闭报名失败")
		return
	}
	response.Success(c, nil)
}

func (h *EventHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("", middleware.AuthRequired())
	{
		auth.POST("/tournaments/:tournament_id/events", middleware.RoleRequired("committee"), h.Create)
		auth.GET("/tournaments/:tournament_id/events", h.List)
		auth.POST("/events/:event_id/close-registration", middleware.RoleRequired("committee"), h.CloseRegistration)
	}
}
