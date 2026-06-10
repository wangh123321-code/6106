package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tt-tournament/backend/internal/middleware"
	"github.com/tt-tournament/backend/internal/service"
	"github.com/tt-tournament/backend/pkg/response"
)

type RegistrationHandler struct {
	svc *service.RegistrationService
}

func NewRegistrationHandler(svc *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{svc: svc}
}

func (h *RegistrationHandler) Register(c *gin.Context) {
	eventID := c.Param("event_id")
	userID := c.GetString("user_id")

	var req struct {
		PartnerID        string `json:"partner_id"`
		QualificationURL string `json:"qualification_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	reg, err := h.svc.Register(c.Request.Context(), eventID, userID, req.PartnerID, req.QualificationURL)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}
	response.Created(c, reg)
}

func (h *RegistrationHandler) BatchApprove(c *gin.Context) {
	eventID := c.Param("event_id")
	var req struct {
		RegistrationIDs []string `json:"registration_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	reviewerID := c.GetString("user_id")
	count, err := h.svc.BatchApprove(c.Request.Context(), eventID, reviewerID, req.RegistrationIDs)
	if err != nil {
		response.InternalError(c, "批量审核失败")
		return
	}
	response.Success(c, gin.H{"approved_count": count})
}

func (h *RegistrationHandler) Reject(c *gin.Context) {
	regID := c.Param("id")
	reviewerID := c.GetString("user_id")
	err := h.svc.Reject(c.Request.Context(), regID, reviewerID)
	if err != nil {
		response.InternalError(c, "拒绝报名失败")
		return
	}
	response.Success(c, nil)
}

func (h *RegistrationHandler) ListByEvent(c *gin.Context) {
	eventID := c.Param("event_id")
	list, err := h.svc.ListByEvent(c.Request.Context(), eventID)
	if err != nil {
		response.InternalError(c, "获取报名列表失败")
		return
	}
	response.Success(c, list)
}

func (h *RegistrationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("", middleware.AuthRequired())
	{
		auth.POST("/events/:event_id/register", h.Register)
		auth.POST("/events/:event_id/batch-approve", middleware.RoleRequired("committee"), h.BatchApprove)
		auth.PUT("/registrations/:id/reject", middleware.RoleRequired("committee"), h.Reject)
		auth.GET("/events/:event_id/registrations", h.ListByEvent)
	}
}
