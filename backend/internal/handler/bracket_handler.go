package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tt-tournament/backend/internal/middleware"
	"github.com/tt-tournament/backend/internal/service"
	"github.com/tt-tournament/backend/pkg/response"
)

type BracketHandler struct {
	svc *service.BracketService
}

func NewBracketHandler(svc *service.BracketService) *BracketHandler {
	return &BracketHandler{svc: svc}
}

func (h *BracketHandler) Generate(c *gin.Context) {
	eventID := c.Param("event_id")
	bracket, err := h.svc.Generate(c.Request.Context(), eventID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, bracket)
}

func (h *BracketHandler) GetByEvent(c *gin.Context) {
	eventID := c.Param("event_id")
	bracket, matches, err := h.svc.GetByEvent(c.Request.Context(), eventID)
	if err != nil {
		response.NotFound(c, "对阵表不存在")
		return
	}
	response.Success(c, gin.H{"bracket": bracket, "matches": matches})
}

func (h *BracketHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("", middleware.AuthRequired())
	{
		auth.POST("/events/:event_id/draw", middleware.RoleRequired("committee"), h.Generate)
		auth.GET("/events/:event_id/bracket", h.GetByEvent)
	}
}
