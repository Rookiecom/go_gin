package router

import (
	"gowork/pkg/controller"
	myjwt "gowork/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func ConfigRouter(router *gin.Engine) {
	user_controller := controller.NewUserController()
	adviser_controller := controller.NewAdviserController()
	order_controller := controller.NewOrderController()
	user := router.Group("/user", myjwt.JWTAuth(1))
	{
		user.POST("/register", user_controller.UserRegister)
		user.POST("/login", user_controller.UserLogin)
		user.GET("/home_page", user_controller.UserHomePage)
		user.POST("/information_edit", user_controller.UserInformationEdit)
		user.GET("/adviser_list", user_controller.GetAdviserList)
		user.POST("/visit_adviser", user_controller.VisitAdviser)
		user.POST("/order_host", user_controller.HostOrder)
		user.POST("/order_finish", user_controller.FinishOrder)

		user.POST("/order_comment", user_controller.CommentOrder)
		user.POST("/order_reward", user_controller.RewardOrder)
		user.POST("/collect_adviser", user_controller.CollectAdviser)
		user.GET("/get_collect_list", user_controller.GetCollectList)
		user.POST("/get_coin_flow", user_controller.UserGetFlow)
		user.POST("/comment_adviser", user_controller.CommentAdviser)
	}
	adviser := router.Group("/adviser", myjwt.JWTAuth(2))
	{
		adviser.POST("/register", adviser_controller.AdviserRegister)
		adviser.POST("/login", adviser_controller.AdviserLogin)
		adviser.GET("/home_page", adviser_controller.AdviserHomePage)
		adviser.POST("/information_edit", adviser_controller.AdviserInformationEdit)
		adviser.POST("/order_take", adviser_controller.TakeOrder)
		adviser.POST("/get_coin_flow", adviser_controller.AdviserGetFlow)
	}
	order := router.Group("/order", myjwt.JWTAuth(0))
	{
		order.GET("/list", order_controller.OrderList)
		order.POST("/single", order_controller.SingleOrder)
	}
}
