package domain

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type IdTable struct {
	Uuid        string `json:"uuid" dynamodbav:"uuid"`
	Type        string `json:"type" dynamodbav:"type"`
	PhoneNumber string `json:"phone_number" dynamodbav:"phone_number"`
}
type FlowOutput struct {
	Time        string `json:"time" dynamodbav:"time"`
	Change      int    `json:"change" dynamodbav:"change"`
	Description string `json:"description" dynamodbav:"description"`
}
type Flow struct {
	Phone_number string `json:"phone_number" dynamodbav:"phone_number"`
	Time         string `json:"time" dynamodbav:"time"`
	Change       int    `json:"change" dynamodbav:"change"`
	Description  string `json:"description" dynamodbav:"description"`
}
type Pagination struct {
	StartHashKey string `json:"start_hash_key" dynamodbav:"start_hash_key"`
	StartTime    string `json:"start_time" dynamodbav:"start_time"`
	PageLimit    int    `json:"page_limit" dynamodbav:"page_limit"`
	SortWay      bool   `json:"sort_way" dynamodbav:"sort_way"`
}

type CommentDisplay struct {
	Commentor string `json:"commentor" dynamodbav:"commentor"`
	Time      string `json:"time" dynamodbav:"time"`
	Score     string `json:"score" dynamodbav:"score"`
	Content   string `json:"content" dynaomdodbav:"content"`
}

var Sess *session.Session
var Svc *dynamodb.DynamoDB

func Create_dynamodb_client() {
	// 创建一个新 session
	Sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
		//Credentials: credentials.NewStaticCredentials("", "", ""),
	})
	// 创建一个新的 DynamoDB 客户端
	if err != nil {
		fmt.Println(err.Error())
	}
	Svc = dynamodb.New(Sess)
}
