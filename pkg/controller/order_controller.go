package controller

import (
	domain "gowork/pkg/domain"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
)

type order_controller struct {
}

func NewOrderController() *order_controller {
	return &order_controller{}
}
func (octl *order_controller) OrderList(c *gin.Context) {
	input := &dynamodb.ScanInput{
		TableName:            aws.String("order_list"),
		ProjectionExpression: aws.String("#id, #name"),
		ExpressionAttributeNames: map[string]*string{
			"#id":   aws.String("id"),
			"#name": aws.String("name"),
		},
	}
	svc := domain.Svc
	result, _ := svc.Scan(input)
	for _, item := range result.Items {
		c.String(http.StatusOK, "name:"+*item["name"].S+"  id:"+*item["id"].S+"\n")
	}
}
func (octl *order_controller) SingleOrder(c *gin.Context) {
	orderid := c.Query("id")
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(orderid)},
		},
		TableName: aws.String("order_list"),
	}
	svc := domain.Svc
	output, _ := svc.GetItem(input)
	var data domain.OrderDisplay
	dynamodbattribute.UnmarshalMap(output.Item, &data)
	c.JSON(http.StatusOK, data)
}
