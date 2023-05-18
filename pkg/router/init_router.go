package router

import (
	"gowork/pkg/controller"
	myjwt "gowork/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func UserConfigRouter(router *gin.Engine) {
	user_controller := controller.NewUserController()
	router.Use(myjwt.JWTAuth(1))
	router.POST("/user_register", user_controller.UserRegister)
	router.POST("/user_login", user_controller.UserLogin)
	router.GET("/user_home_page", user_controller.UserHomePage)
	router.POST("/user_information_edit", user_controller.UserInformationEdit)
	router.GET("/adviser_list", user_controller.GetAdviserList)
	router.POST("/visit_adviser", user_controller.VisitAdviser)
}
func AdviserConfigRouter(router *gin.Engine) {
	adviser_controller := controller.NewAdviserController()
	router.Use(myjwt.JWTAuth(2))
	router.POST("/adviser_register", adviser_controller.AdviserRegister)
	router.POST("/adviser_login", adviser_controller.AdviserLogin)
	router.GET("/adviser_home_page", adviser_controller.AdviserHomePage)
	router.POST("adviser_information_edit", adviser_controller.AdviserInformationEdit)
}

/*
func Pubilic(router *gin.Engine) {
	router.Use(myjwt.JWTAuth(0))
	user_controller := controller.NewUserController()
	adviser_controller := controller.NewAdviserController()
	router.GET("/user_home_page", user_controller.UserHomePage)
	router.GET("/adviser_home_page", adviser_controller.AdviserHomePage)
}
*/
