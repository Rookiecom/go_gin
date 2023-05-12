package router

import (
	"fmt"
	"gowork/pkg/controller"
	myjwt "gowork/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func UserConfigRouter(router *gin.Engine) {
	user_controller := controller.NewUserController()
	router.Use(myjwt.JWTAuth(1))
	router.POST("/user_register", user_controller.UserRegister)
	router.POST("/user_login", user_controller.UserLogin)
	router.POST("/user_home_page", user_controller.UserHomePage)
}
func AdviserConfigRouter(router *gin.Engine) {
	adviser_controller := controller.NewAdviserController()
	fmt.Println("adfsasf")
	router.Use(myjwt.JWTAuth(2))
	router.POST("/adviser_register", adviser_controller.AdviserRegister)
	router.POST("/adviser_login", adviser_controller.AdviserLogin)
	router.POST("/adviser_home_page", adviser_controller.AdviserHomePage)
}
