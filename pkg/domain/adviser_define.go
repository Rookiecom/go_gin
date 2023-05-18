package domain

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type AdviserEdit struct {
	Name      string `json:"name" dynamodbbav:"name"`
	WorkState bool   `json:"workstate" dynamodbbav:"workstate"`
}

type AdviserDisEdit struct {
	Coin       int    `json:"coin" dynamodbbav:"coin"`
	OrderNum   int    `json:"ordernum" dynamodbbav:"ordernum"`
	Score      int    `json:"score" dynamodbbav:"score"`
	CommentNum int    `json:"commentnum" dynamodbbav:"commentnum"`
	Uuid       string `json:"uuid" dynamodbav:"uuid"`
}
type AdviserRequest struct {
	Loginreq
	AdviserEdit
}
type AdviserRegister struct {
	AdviserRequest
	AdviserDisEdit
}
type AdviserHomePage struct {
	AdviserEdit
	AdviserDisEdit
}
type AdviserRequestUpdate struct {
	IsEdit bool `json:"is_edit" dynamodbbav:"is_edit"`
	AdviserEdit
}
func (aru *AdviserRequestUpdate) Update(phone_number string) (*dynamodb.UpdateItemOutput, error) {
	update_item_input := &dynamodb.UpdateItemInput{
		TableName: aws.String("adviser"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(phone_number)},
		},
		UpdateExpression: aws.String("SET #na=:name,#wo=:workstate"),
		ExpressionAttributeNames: map[string]*string{
			"#na": aws.String("name"),
			"#wo": aws.String("workstate"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(aru.Name),
			},
			":workstate": {
				BOOL: aws.Bool(aru.WorkState),
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	return Svc.UpdateItem(update_item_input)

}
