package db

import (
	"fmt"
	"gowork/pkg/domain"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

func OrderRegister(order *domain.Order, phone_number string) error {
	order.Status = "nobody take"
	uuid, _ := uuid.NewRandom()
	order.ID = uuid.String()
	order.OrderHostPhone = phone_number
	order_item, _ := dynamodbattribute.MarshalMap(order)
	order_input := &dynamodb.PutItemInput{
		Item:      order_item,
		TableName: aws.String("order_list"),
	}
	svc := domain.Svc
	_, ok1 := svc.PutItem(order_input)
	if ok1 != nil {
		return ok1
	}
	help_order_query := &domain.AuxiliaryQuery{
		PhoneNumber: phone_number,
		Name:        order.Name,
		ID:          order.ID,
		Type:        "user",
	}
	help__item, _ := dynamodbattribute.MarshalMap(help_order_query)
	help_input := &dynamodb.PutItemInput{
		Item:      help__item,
		TableName: aws.String("query_order"),
	}
	_, ok2 := svc.PutItem(help_input)
	if ok2 != nil {
		return ok2
	}
	return nil
}

func OrderTake(id string, name string, phone_number string) error {
	svc := domain.Svc
	update_order_input := &dynamodb.UpdateItemInput{
		TableName: aws.String("order_list"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		UpdateExpression: aws.String("SET #st=:status,#oth=:val"),
		ExpressionAttributeNames: map[string]*string{
			"#st":  aws.String("status"),
			"#oth": aws.String("order_taker_phone"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {S: aws.String(name + " take order")},
			":val":    {S: aws.String(phone_number)},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	update_adviser_input := &dynamodb.UpdateItemInput{
		TableName: aws.String("adviser_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(phone_number)},
		},
		UpdateExpression: aws.String("SET ordernum=ordernum+:val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {N: aws.String("1")},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	_, ok1 := svc.UpdateItem(update_order_input)
	if ok1 != nil {
		return ok1
	}
	_, ok2 := svc.UpdateItem(update_adviser_input)
	if ok2 != nil {
		return ok2
	}
	query_data := domain.AuxiliaryQuery{
		PhoneNumber: phone_number,
		Name:        name,
		ID:          id,
		Type:        "adviser",
	}
	query_input, _ := dynamodbattribute.MarshalMap(query_data)
	query_item := &dynamodb.PutItemInput{
		TableName: aws.String("query_order"),
		Item:      query_input,
	}
	svc.PutItem(query_item)
	return nil
}
func OrderFinish(orderid string) error {
	svc := domain.Svc
	update_order_input := &dynamodb.UpdateItemInput{
		TableName: aws.String("order_list"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(orderid)},
		},
		UpdateExpression: aws.String("SET #st=:status"),
		ExpressionAttributeNames: map[string]*string{
			"#st": aws.String("status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {S: aws.String("finish")},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	_, ok1 := svc.UpdateItem(update_order_input)
	if ok1 != nil {
		return ok1
	}
	get_order_input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(orderid)},
		},
		TableName: aws.String("order_list"),
	}
	orderoutput, ok2 := svc.GetItem(get_order_input)
	if ok2 != nil {
		return ok2
	}
	var orderdata domain.Order
	dynamodbattribute.UnmarshalMap(orderoutput.Item, &orderdata)
	user_expression := fmt.Sprintf("SET %s=%s- :val", "coin", "coin")
	update_user_input := &dynamodb.UpdateItemInput{
		TableName: aws.String("user_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(orderdata.OrderHostPhone)},
		},
		UpdateExpression: aws.String(user_expression),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {N: aws.String(fmt.Sprintf("%d", orderdata.Price))},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	_, ok3 := svc.UpdateItem(update_user_input)
	if ok3 != nil {
		return ok3
	}
	user_flow := &domain.Flow{
		Phone_number: orderdata.OrderHostPhone,
		Time:         time.Now().Format("2006-01-02T15:04:05Z"),
		Change:       -orderdata.Price,
		Description:  "pay order " + orderid,
	}
	user_flow_item, _ := dynamodbattribute.MarshalMap(user_flow)
	put_user_input := &dynamodb.PutItemInput{
		TableName: aws.String("coin_flow"),
		Item:      user_flow_item,
	}
	_, ok4 := svc.PutItem(put_user_input)
	if ok4 != nil {
		return ok4
	}
	update_adviser_input := &dynamodb.UpdateItemInput{
		TableName: aws.String("adviser_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(orderdata.OrderTakerPhone)},
		},
		UpdateExpression: aws.String("SET coin = coin + :val1 , order_finish_num=order_finish_num+:val2"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val1": {N: aws.String(fmt.Sprintf("%d", orderdata.Price))},
			":val2": {N: aws.String("1")},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	_, ok5 := svc.UpdateItem(update_adviser_input)
	if ok5 != nil {
		fmt.Println(ok5.Error())
		return ok5
	}
	adviser_flow := &domain.Flow{
		Phone_number: orderdata.OrderTakerPhone,
		Time:         time.Now().Format("2006-01-02T15:04:05Z"),
		Change:       orderdata.Price,
		Description:  "finish order " + orderid,
	}
	adviser_flow_item, _ := dynamodbattribute.MarshalMap(adviser_flow)
	put_adviser_input := &dynamodb.PutItemInput{
		TableName: aws.String("coin_flow"),
		Item:      adviser_flow_item,
	}
	_, ok6 := svc.PutItem(put_adviser_input)

	if ok6 != nil {
		//fmt.Println(ok6.Error())
		return ok6
	}
	return nil
}
func OrderComment(comment *domain.CommentRequest) error {
	svc := domain.Svc
	updata_item := &dynamodb.UpdateItemInput{
		TableName: aws.String("order_list"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(comment.ID)},
		},
		UpdateExpression: aws.String("SET comment_text=:val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {S: aws.String(comment.CommentText)},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	_, ok := svc.UpdateItem(updata_item)
	return ok
}
func OrderReward(reward *domain.RewardRequest) error {
	svc := domain.Svc
	query_order_input := &dynamodb.GetItemInput{
		TableName: aws.String("order_list"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(reward.ID)},
		},
	}
	query_order, ok1 := svc.GetItem(query_order_input)
	if ok1 != nil {
		return ok1
	}
	user_updata_item := &dynamodb.UpdateItemInput{
		TableName: aws.String("user_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(*query_order.Item["order_host_phone"].S)},
		},
		UpdateExpression: aws.String("SET coin=coin - :val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {N: aws.String(fmt.Sprintf("%d", reward.Reward))},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	_, ok2 := svc.UpdateItem(user_updata_item)
	if ok2 != nil {
		return ok2
	}
	user_flow := &domain.Flow{
		Phone_number: *query_order.Item["order_host_phone"].S,
		Time:         time.Now().Format("2006-01-02T15:04:05Z"),
		Change:       -reward.Reward,
		Description:  "reward order " + reward.ID,
	}
	user_flow_item, _ := dynamodbattribute.MarshalMap(user_flow)
	put_user_input := &dynamodb.PutItemInput{
		TableName: aws.String("coin_flow"),
		Item:      user_flow_item,
	}
	_, ok3 := svc.PutItem(put_user_input)
	if ok3 != nil {
		return ok3
	}
	adviser_updata_item := &dynamodb.UpdateItemInput{
		TableName: aws.String("adviser_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(*query_order.Item["order_taker_phone"].S)},
		},
		UpdateExpression: aws.String("SET coin=coin + :val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {N: aws.String(fmt.Sprintf("%d", reward.Reward))},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	_, ok4 := svc.UpdateItem(adviser_updata_item)
	if ok4 != nil {
		return ok4
	}
	adviser_flow := &domain.Flow{
		Phone_number: *query_order.Item["order_taker_phone"].S,
		Time:         time.Now().Format("2006-01-02T15:04:05Z"),
		Change:       reward.Reward,
		Description:  "rewarded for order " + reward.ID,
	}
	adviser_flow_item, _ := dynamodbattribute.MarshalMap(adviser_flow)
	put_adviser_input := &dynamodb.PutItemInput{
		TableName: aws.String("coin_flow"),
		Item:      adviser_flow_item,
	}
	_, ok5 := svc.PutItem(put_adviser_input)
	if ok5 != nil {
		return ok5
	}
	return nil
}
