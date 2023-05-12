package main

import (
	"gowork/pkg/domain"
	myrouter "gowork/pkg/router"

	"github.com/gin-gonic/gin"
)

func main() {
	domain.Create_dynamodb_client()
	user_router := gin.Default()
	adviser_router := gin.Default()
	myrouter.UserConfigRouter(user_router)
	myrouter.AdviserConfigRouter(adviser_router)
	go func() {
		user_router.Run(":8080")
	}()
	go func() {
		adviser_router.Run(":8090")
	}()
	select {}
}
