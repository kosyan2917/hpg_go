package handlers

import (
	"github.com/gin-gonic/gin"
	controller "hpg_backend_go/controllers"
	"hpg_backend_go/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}

func UnAuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("board", middlewares.BoardAuth(), controller.Board())
}

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("users/me", middlewares.UserAuthenticate(), controller.Me())
}
