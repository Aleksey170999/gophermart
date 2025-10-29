package handler

import (
	"github.com/Aleksey170999/go-loyaty/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Handler struct {
	services   *service.Service
	jwtService *service.JWTService
	logger     *zap.Logger
}

func NewHandler(services *service.Service, jwtService *service.JWTService, logger *zap.Logger) *Handler {
	return &Handler{
		services:   services,
		jwtService: jwtService,
		logger:     logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	jwtMiddleware := NewJWTMiddleware(h.jwtService)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/api/user")
	{
		auth.POST("/register", h.SignUp)
		auth.POST("/login", h.SignIn)
	}

	authorized := router.Group("/api")
	authorized.Use(jwtMiddleware.Auth())
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
