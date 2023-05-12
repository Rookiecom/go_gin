package domain

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Loginreq struct {
	PhoneNumber string `json:"phone_number" dynamodbbav:"phone_number" `
	Password    string `json:"password" dynamodbbav:"password"`
}

type UserEdit struct {
	Name   string `json:"name" dynamodbbav:"name"`
	Birth  string `json:"birth" dynamodbbav:"birth"`
	Gender string `json:"gender" dynamodbbav:"gender"`
	Bio    string `json:"bio" dynamodbbav:"bio"`
	About  string `json:"about" dynamodbbav:"about"`
}
type UserDisEdit struct {
	Coin int `json:"coin" dynamodbbav:"coin"`
}
type UserRequest struct {
	Loginreq
	UserEdit
}
type UserRegister struct {
	UserRequest
	UserDisEdit
}
type UserHomePage struct {
	UserEdit
	UserDisEdit
}
type UserRequestUpdate struct {
	IsEdit bool `json:"is_edit" dynamodbbav:"is_edit"`
	UserEdit
}

func UserDisEditInit() *UserDisEdit {
	return &UserDisEdit{
		Coin: 0,
	}
}
func (uru *UserRequestUpdate) Update(phone_number string) (*dynamodb.UpdateItemOutput, error) {
	update_item, _ := dynamodbattribute.MarshalMap(&uru.UserEdit)
	update_item_input := &dynamodb.UpdateItemInput{
		TableName: aws.String("user"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(phone_number)},
		},
		UpdateExpression: aws.String("SET #na=:name,#bir=:birth,#gen=:gender,#bio=:bio,#ab=:about"),
		ExpressionAttributeNames: map[string]*string{
			"#na":  aws.String("name"),
			"#bir": aws.String("birth"),
			"#gen": aws.String("gender"),
			"#bio": aws.String("bio"),
			"#ab":  aws.String("about"),
		},
		ExpressionAttributeValues: update_item,
		ReturnValues:              aws.String("ALL_NEW"),
	}
	return Svc.UpdateItem(update_item_input)

}
