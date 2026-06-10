package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tt-tournament/backend/internal/middleware"
	"github.com/tt-tournament/backend/internal/service"
	"github.com/tt-tournament/backend/pkg/response"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required,min=6"`
		DisplayName string `json:"display_name" binding:"required"`
		Role        string `json:"role" binding:"required,oneof=player referee committee"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, err := h.svc.Register(c.Request.Context(), req.Username, req.Password, req.DisplayName, req.Role)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}
	response.Created(c, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	token, user, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.Success(c, gin.H{"token": token, "user": user})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.svc.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) UpdateRanking(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Ranking int `json:"ranking" binding:"min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	err := h.svc.UpdateRanking(c.Request.Context(), id, req.Ranking)
	if err != nil {
		response.InternalError(c, "更新排名失败")
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	role := c.Query("role")
	users, err := h.svc.ListByRole(c.Request.Context(), role)
	if err != nil {
		response.InternalError(c, "获取用户列表失败")
		return
	}
	response.Success(c, users)
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)

	auth := rg.Group("", middleware.AuthRequired())
	{
		auth.GET("/profile", h.GetProfile)
		auth.GET("/users", middleware.RoleRequired("committee"), h.ListUsers)
		auth.PUT("/users/:id/ranking", middleware.RoleRequired("committee"), h.UpdateRanking)
	}
}
