package db

import (
	"gowork/pkg/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var TableStatus map[string]bool

func TableStatusInit() {
	TableStatus = make(map[string]bool)
	const tablenum = 8
	tablename := [tablenum]string{
		"user_info",
		"adviser_info",
		"all_id",
		"order_list",
		"query_order",
		"user_collect_adviser",
		"coin_flow",
		"adviser_comment",
	}
	var err [tablenum]error
	for i := 0; i < tablenum; i++ {
		err[i] = TableDescribe(tablename[i])
		if err[i] == nil {
			TableStatus[tablename[i]] = true
		} else {
			TableStatus[tablename[i]] = false
		}
	}

}

//本文件是复杂的数据库访问的集成
func CreateUserTable() error {
	if TableStatus["user_info"] {
		return nil
	}
	svc := domain.Svc
	user_create := &dynamodb.CreateTableInput{
		TableName: aws.String("user_info"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("phone_number"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("phone_number"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	_, err := svc.CreateTable(user_create)
	if err == nil {
		TableStatus["user_info"] = true
	}
	return err
}
func CreateAdviserTable() error {
	if TableStatus["adviser_info"] {
		return nil
	}
	svc := domain.Svc
	adviser_create := &dynamodb.CreateTableInput{
		TableName: aws.String("adviser_info"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("phone_number"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("phone_number"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	_, err := svc.CreateTable(adviser_create)
	if err == nil {
		TableStatus["adviser_info"] = true
	}
	return err
}
func CreateIdTable() error {
	if TableStatus["all_id"] {
		return nil
	}
	svc := domain.Svc
	id_create := &dynamodb.CreateTableInput{
		TableName: aws.String("all_id"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("uuid"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("uuid"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	_, err := svc.CreateTable(id_create)
	if err == nil {
		TableStatus["all_id"] = true
	}
	return err
}
func CreateOrderTable() error {
	svc := domain.Svc
	if TableStatus["order_list"] {
		return nil
	}
	order_table_create := &dynamodb.CreateTableInput{
		TableName: aws.String("order_list"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	_, err := svc.CreateTable(order_table_create)
	if err == nil {
		TableStatus["order_list"] = true
	}
	return err
}
func CreateSelectOrder() error {
	if TableStatus["query_order"] {
		return nil
	}
	svc := domain.Svc
	adviser_create := &dynamodb.CreateTableInput{
		TableName: aws.String("query_order"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("phone_number"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("type"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("phone_number"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("type"),
				KeyType:       aws.String("range"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	_, err := svc.CreateTable(adviser_create)
	if err == nil {
		TableStatus["query_order"] = true
	}
	return err
}
func CreateCollect() error {
	if TableStatus["user_collect_adviser"] {
		return nil
	}
	svc := domain.Svc
	create_input := &dynamodb.CreateTableInput{
		TableName: aws.String("user_collect_adviser"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("phone_number"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("phone_number"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	_, err := svc.CreateTable(create_input)
	return err
}
func CreateFlowTable() error {
	if TableStatus["coin_flow"] {
		return nil
	}
	svc := domain.Svc
	create_input := &dynamodb.CreateTableInput{
		TableName: aws.String("coin_flow"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("phone_number"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("time"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("phone_number"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("time"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	_, err := svc.CreateTable(create_input)
	return err
}
func CreateAdviserCommentTable() error {
	if TableStatus["adviser_comment"] {
		return nil
	}
	svc := domain.Svc
	create_input := &dynamodb.CreateTableInput{
		TableName: aws.String("adviser_comment"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("phone_number"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("time"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("phone_number"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("time"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	_, err := svc.CreateTable(create_input)
	return err
}
func TableDescribe(tablename string) error {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tablename),
	}
	svc := domain.Svc
	_, err := svc.DescribeTable(input)
	return err
}
