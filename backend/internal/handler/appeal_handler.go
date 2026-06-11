package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tt-tournament/backend/internal/middleware"
	"github.com/tt-tournament/backend/internal/service"
	"github.com/tt-tournament/backend/pkg/response"
)

type AppealHandler struct {
	svc *service.AppealService
}

func NewAppealHandler(svc *service.AppealService) *AppealHandler {
	return &AppealHandler{svc: svc}
}

func (h *AppealHandler) Create(c *gin.Context) {
	appellantID := c.GetString("user_id")

	var req service.CreateAppealInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	appeal, err := h.svc.Create(c.Request.Context(), appellantID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, appeal)
}

func (h *AppealHandler) List(c *gin.Context) {
	filter := map[string]string{
		"match_id":     c.Query("match_id"),
		"appellant_id": c.Query("appellant_id"),
		"status":       c.Query("status"),
	}

	role := c.GetString("role")
	userID := c.GetString("user_id")
	if role == "player" {
		filter["appellant_id"] = userID
	}

	list, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, "获取申诉列表失败")
		return
	}
	response.Success(c, list)
}

func (h *AppealHandler) Review(c *gin.Context) {
	appealID := c.Param("id")
	reviewerID := c.GetString("user_id")
	reviewerName := c.GetString("username")

	var req service.ReviewAppealInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	appeal, err := h.svc.Review(c.Request.Context(), appealID, reviewerID, reviewerName, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, appeal)
}

func (h *AppealHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("", middleware.AuthRequired())
	{
		auth.POST("/appeals", middleware.RoleRequired("player", "committee"), h.Create)
		auth.GET("/appeals", h.List)
		auth.PUT("/appeals/:id/review", middleware.RoleRequired("committee"), h.Review)
	}
}
