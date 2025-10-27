package handler

import (
	"github.com/Aleksey170999/go-loyaty/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/api/user")
	{
		auth.POST("/register", h.SignUp)
		auth.POST("/login", h.SignIn)
	}

	authorized := router.Group("/api", JWTAuthMiddleware())
	{
		user := authorized.Group("/user")
		{
			user.GET("orders/:number", h.GetOrderInfo)
			user.POST("orders/", h.CreateOrder)
			user.GET("orders/", h.GetOrdersList)
			user.GET("/balance", h.GetUserBalance)
			user.POST("/balance/withdraw", h.ProcessWithdrawal)
			user.GET("/withdrawals", h.GetWithdrawals)
		}
	}

	return router
}
