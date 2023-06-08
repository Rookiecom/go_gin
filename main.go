package main

import (
	mydb "gowork/pkg/db"
	"gowork/pkg/domain"
	myrouter "gowork/pkg/router"

	"github.com/gin-gonic/gin"
)

func main() {
	domain.Create_dynamodb_client()
	mydb.TableStatusInit()
	router := gin.Default()
	myrouter.ConfigRouter(router)
	router.Run(":8080")
}
