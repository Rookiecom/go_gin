package domain

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var Sess *session.Session
var Svc *dynamodb.DynamoDB

func Create_dynamodb_client() {
	// 创建一个新 session
	fmt.Print("sdaf")
	Sess, _ = session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
		//Credentials: credentials.NewStaticCredentials("", "", ""),
	})
	// 创建一个新的 DynamoDB 客户端
	Svc = dynamodb.New(Sess)
}
