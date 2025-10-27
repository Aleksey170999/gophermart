package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	var input string

	if err := c.BindPlain(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userID := c.GetInt("user_id")

	if userID == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	_, err := h.services.Orders.CreateOrder(input, userID)
	if err != nil {
		switch err.Error() {
		case "invalid order number format":
			newErrorResponse(c, http.StatusUnprocessableEntity, err.Error()) // 422
		case "order already uploaded by this user":
			c.Status(http.StatusOK) // 200
		case "order already uploaded by another user":
			newErrorResponse(c, http.StatusConflict, err.Error()) // 409
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error()) // 500
		}
		return
	}

	c.Status(http.StatusAccepted) // 202
}

func (h *Handler) GetUserBalance(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get current balance
	current, err := h.services.Orders.GetUserBalance(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get balance")
		return
	}

	// Get withdrawn sum
	withdrawn, err := h.services.Withdrawals.GetWithdrawnSum(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to get withdrawn sum")
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
		newErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	orders, err := h.services.Orders.GetOrdersList(userID)
	if err != nil {
		switch err.Error() {
		case "invalid order number format":
			newErrorResponse(c, http.StatusUnprocessableEntity, err.Error()) // 422
		case "order already uploaded by this user":
			c.Status(http.StatusOK) // 200
		case "order already uploaded by another user":
			newErrorResponse(c, http.StatusConflict, err.Error()) // 409
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error()) // 500
		}
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
		newErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	orderNumber := c.Param("number")
	if orderNumber == "" {
		newErrorResponse(c, http.StatusBadRequest, "Order number is required")
		return
	}
	order, err := h.services.Orders.GetOrderInfo(orderNumber)
	if err != nil {
		switch err.Error() {
		case "invalid order number format":
			newErrorResponse(c, http.StatusUnprocessableEntity, err.Error()) // 422
		case "order already uploaded by this user":
			c.Status(http.StatusOK) // 200
		case "order already uploaded by another user":
			newErrorResponse(c, http.StatusConflict, err.Error()) // 409
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error()) // 500
		}
		return
	}
	c.JSON(http.StatusOK, order)
}
