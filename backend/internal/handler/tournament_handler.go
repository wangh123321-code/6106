package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tt-tournament/backend/internal/middleware"
	"github.com/tt-tournament/backend/internal/service"
	"github.com/tt-tournament/backend/pkg/response"
)

type TournamentHandler struct {
	svc *service.TournamentService
}

func NewTournamentHandler(svc *service.TournamentService) *TournamentHandler {
	return &TournamentHandler{svc: svc}
}

func (h *TournamentHandler) Create(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Location  string `json:"location" binding:"required"`
		StartDate string `json:"start_date" binding:"required"`
		EndDate   string `json:"end_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response.BadRequest(c, "开始日期格式错误")
		return
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response.BadRequest(c, "结束日期格式错误")
		return
	}

	userID := c.GetString("user_id")
	t, err := h.svc.Create(c.Request.Context(), req.Name, req.Location, start, end, userID)
	if err != nil {
		response.InternalError(c, "创建赛事失败")
		return
	}
	response.Created(c, t)
}

func (h *TournamentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	t, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "赛事不存在")
		return
	}
	response.Success(c, t)
}

func (h *TournamentHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.InternalError(c, "获取赛事列表失败")
		return
	}
	response.Success(c, list)
}

func (h *TournamentHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required,oneof=draft open registration_closed drawn in_progress completed published"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		response.InternalError(c, "更新状态失败")
		return
	}
	response.Success(c, nil)
}

func (h *TournamentHandler) Publish(c *gin.Context) {
	id := c.Param("id")
	err := h.svc.Publish(c.Request.Context(), id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *TournamentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("", middleware.AuthRequired())
	{
		auth.POST("/tournaments", middleware.RoleRequired("committee"), h.Create)
		auth.GET("/tournaments", h.List)
		auth.GET("/tournaments/:id", h.GetByID)
		auth.PUT("/tournaments/:id/status", middleware.RoleRequired("committee"), h.UpdateStatus)
		auth.POST("/tournaments/:id/publish", middleware.RoleRequired("committee"), h.Publish)
	}
}
