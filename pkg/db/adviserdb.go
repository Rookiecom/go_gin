package db

import (
	"gowork/pkg/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

func AdviserRegisterDB(adviserreq *domain.AdviserRequest) error {
	svc := domain.Svc
	newuuid, _ := uuid.NewRandom()
	adviser_register := domain.AdviserRegister{
		AdviserRequest: *adviserreq,
		AdviserDisEdit: domain.AdviserDisEdit{
			Coin:           0,
			OrderNum:       0,
			OrderFinishNum: 0,
			Score:          0,
			CommentNum:     0,
			Uuid:           newuuid.String(),
		},
	}
	adviser_item, _ := dynamodbattribute.MarshalMap(adviser_register)
	adviser_item_input := &dynamodb.PutItemInput{
		Item:      adviser_item,
		TableName: aws.String("adviser_info"),
	}
	_, err1 := svc.PutItem(adviser_item_input)
	if err1 != nil {
		return err1
	}
	id_table := domain.IdTable{
		Uuid:        newuuid.String(),
		Type:        "adviser",
		PhoneNumber: adviserreq.PhoneNumber,
	}
	id_table_item, _ := dynamodbattribute.MarshalMap(id_table)
	id_table_input := &dynamodb.PutItemInput{
		Item:      id_table_item,
		TableName: aws.String("all_id"),
	}
	_, err2 := svc.PutItem(id_table_input)
	if err2 != nil {
		return err2
	}
	return nil
}
func AdviserUpdate(aru *domain.AdviserRequestUpdate, phone_number string) (*dynamodb.UpdateItemOutput, error) {
	svc := domain.Svc
	update_item_input := &dynamodb.UpdateItemInput{
		TableName: aws.String("adviser_update"),
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
	return svc.UpdateItem(update_item_input)

}
