package handlers

import (
	"github.com/gin-gonic/gin"
	controller "hpg_backend_go/controllers"
	"hpg_backend_go/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
	incomingRoutes.GET("users/refresh", controller.Refresh())
}

func UnAuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("board", middlewares.BoardAuth(), controller.Board())
	incomingRoutes.GET("media/:dir/*asset", controller.Media())
	incomingRoutes.GET("users", controller.Users())
	incomingRoutes.GET("field/:id", controller.Field())
}

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("users/me", middlewares.UserAuthenticate(), controller.Me())
	incomingRoutes.PUT("users/change-avatar", middlewares.UserAuthenticate(), controller.ChangeAvatar())
	incomingRoutes.PUT("users/change-color", middlewares.UserAuthenticate(), controller.ChangeColor())
	incomingRoutes.GET("roll", middlewares.UserAuthenticate(), controller.Roll())
	incomingRoutes.PUT("game", middlewares.UserAuthenticate(), controller.SetCurrentGame())
	incomingRoutes.GET("game", middlewares.UserAuthenticate(), controller.GetCurrentGame())
}
