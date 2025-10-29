package handler

import (
	"net/http"

	"github.com/Aleksey170999/go-loyaty/internal/apperror"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) handleAppError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		if appErr.HTTPStatus == http.StatusOK {
			c.Status(appErr.HTTPStatus)
			return
		}
		newErrorResponse(c, appErr.HTTPStatus, appErr.Message)
		return
	}
	h.logger.Error("Internal server error", zap.Error(err))
	newErrorResponse(c, http.StatusInternalServerError, "internal server error")
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var input string

	if err := c.BindPlain(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := c.GetInt("user_id")
	if userID == 0 {
		h.handleAppError(c, apperror.ErrUnauthorized)
		return
	}

	_, err := h.services.Orders.CreateOrder(input, userID)
	if err != nil {
		h.handleAppError(c, err)
		return
	}

	c.Status(http.StatusAccepted)
}

func (h *Handler) GetUserBalance(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		h.handleAppError(c, apperror.ErrUnauthorized)
		return
	}

	// Get current balance
	current, err := h.services.Orders.GetUserBalance(userID)
	if err != nil {
		h.handleAppError(c, apperror.Wrap(err, "failed to get balance", http.StatusInternalServerError))
		return
	}

	// Get withdrawn sum
	withdrawn, err := h.services.Withdrawals.GetWithdrawnSum(userID)
	if err != nil {
		h.handleAppError(c, apperror.Wrap(err, "failed to get withdrawn sum", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"current":   current,
		"withdrawn": withdrawn,
	})
}

func (h *Handler) GetOrdersList(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		h.handleAppError(c, apperror.ErrUnauthorized)
		return
	}

	orders, err := h.services.Orders.GetOrdersList(userID)
	if err != nil {
		h.handleAppError(c, err)
		return
	}

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetOrderInfo(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		h.handleAppError(c, apperror.ErrUnauthorized)
		return
	}

	orderNumber := c.Param("number")
	if orderNumber == "" {
		newErrorResponse(c, http.StatusBadRequest, "order number is required")
		return
	}

	order, err := h.services.Orders.GetOrderInfo(orderNumber)
	if err != nil {
		h.handleAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, order)
}
