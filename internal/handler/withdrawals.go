package handler

import (
	"net/http"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ProcessWithdrawal(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.WithdrawalRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Order == "" || req.Sum <= 0 {
		newErrorResponse(c, http.StatusUnprocessableEntity, "invalid order number or sum")
		return
	}

	err := h.services.Withdrawals.ProcessWithdrawal(userID, req.Order, req.Sum)
	if err != nil {
		switch err.Error() {
		case "insufficient funds":
			newErrorResponse(c, http.StatusPaymentRequired, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetWithdrawals(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	withdrawals, err := h.services.Withdrawals.GetWithdrawals(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}
