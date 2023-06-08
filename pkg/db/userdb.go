package db

import (
	"fmt"
	"gowork/pkg/domain"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

func UserRegisterDB(userreq *domain.UserRequest) error {
	svc := domain.Svc
	newuuid, _ := uuid.NewRandom()
	user_register := domain.UserRegister{
		UserRequest: *userreq,
		UserDisEdit: domain.UserDisEdit{
			Coin: 0,
			Uuid: newuuid.String(),
		},
	}
	useritem, _ := dynamodbattribute.MarshalMap(user_register)
	useriteminput := &dynamodb.PutItemInput{
		Item:      useritem,
		TableName: aws.String("user_info"),
	}
	_, err1 := svc.PutItem(useriteminput)
	if err1 != nil {
		return err1
	}
	id_table := domain.IdTable{
		Uuid:        newuuid.String(),
		Type:        "user",
		PhoneNumber: user_register.PhoneNumber,
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
func UserUpdate(uru *domain.UserRequestUpdate, phone_number string) (*dynamodb.UpdateItemOutput, error) {
	svc := domain.Svc
	update_item_input := &dynamodb.UpdateItemInput{
		TableName: aws.String("user_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(phone_number)},
		},
		UpdateExpression: aws.String("SET #na=:name,#bir=:birth,#gen=:gender,#bio=:bio,#ab=:about"),
		ExpressionAttributeNames: map[string]*string{
			"#ab":  aws.String("about"),
			"#bio": aws.String("bio"),
			"#bir": aws.String("birth"),
			"#gen": aws.String("gender"),
			"#na":  aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(uru.Name),
			},
			":birth": {
				S: aws.String(uru.Birth),
			},
			":bio": {
				S: aws.String(uru.Bio),
			},
			":gender": {
				S: aws.String(uru.Gender),
			},
			":about": {
				S: aws.String(uru.About),
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	return svc.UpdateItem(update_item_input)

}
func CollectAdviser(user_phone_number, adviserid string) error {
	svc := domain.Svc
	get_record_input := &dynamodb.GetItemInput{
		TableName: aws.String("user_collect_adviser"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(user_phone_number)},
		},
	}
	get_record_output, ok1 := svc.GetItem(get_record_input)
	if ok1 != nil {
		return ok1
	}
	if get_record_output.Item == nil { //未收藏过
		type row struct {
			PhoneNumber string   `json:"phone_number" dynamodbav:"phone_number"`
			AdviserList []string `json:"adviser_list" dynamodbav:"adviser_list"`
		} //表中一行的属性
		item := &row{
			PhoneNumber: user_phone_number,
			AdviserList: []string{adviserid},
		}
		item_attr, _ := dynamodbattribute.MarshalMap(item)
		put_item := &dynamodb.PutItemInput{
			TableName: aws.String("user_collect_adviser"),
			Item:      item_attr,
		}
		_, ok2 := svc.PutItem(put_item)
		return ok2
	} else { //收藏过
		update_item := &dynamodb.UpdateItemInput{
			TableName: aws.String("user_collect_adviser"),
			Key: map[string]*dynamodb.AttributeValue{
				"phone_number": {S: aws.String(user_phone_number)},
			},
			UpdateExpression: aws.String("SET #ri = list_append(#ri,:val)"),
			ExpressionAttributeNames: map[string]*string{
				"#ri": aws.String("adviser_list"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":val": {
					L: []*dynamodb.AttributeValue{
						{S: aws.String(adviserid)},
					},
				},
			},
			ReturnValues: aws.String("ALL_NEW"),
		}
		_, ok2 := svc.UpdateItem(update_item)
		return ok2
	}
}
func GetFlow(phone_number string, page_request *domain.Pagination) (*dynamodb.QueryOutput, error) {
	svc := domain.Svc
	var startkey map[string]*dynamodb.AttributeValue
	if page_request.StartHashKey == "" {
		startkey = nil
	} else {
		startkey = map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(page_request.StartHashKey)},
			"time":         {S: aws.String(page_request.StartTime)},
		}
	}
	query_input := &dynamodb.QueryInput{
		TableName: aws.String("coin_flow"),

		KeyConditionExpression: aws.String("phone_number=:val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {S: aws.String(phone_number)},
		},
		Limit:                    aws.Int64(int64(page_request.PageLimit)),
		ExclusiveStartKey:        startkey,
		ProjectionExpression:     aws.String("#t,change,description"),
		ExpressionAttributeNames: map[string]*string{"#t": aws.String("time")},
		ScanIndexForward:         aws.Bool(page_request.SortWay),
	}
	return svc.Query(query_input)
}

func CommentAdviser(commentor string, comment_request *domain.AdviserCommentRequeat) error {
	svc := domain.Svc
	get_input := &dynamodb.GetItemInput{
		TableName: aws.String("all_id"),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {S: aws.String(comment_request.AdviserId)},
		},
	}
	get_output, ok1 := svc.GetItem(get_input)
	if ok1 != nil {
		return ok1
	}
	comment := &domain.AdviserComment{
		PhoneNumber: *get_output.Item["phone_number"].S,
		Content:     comment_request.Content,
		Score:       float64(comment_request.Score),
		Time:        time.Now().Format("2006-01-02T15:04:05Z"),
		Commentor:   commentor,
	}
	item, _ := dynamodbattribute.MarshalMap(comment) //attention
	put_input := &dynamodb.PutItemInput{
		TableName: aws.String("adviser_comment"),
		Item:      item,
	}
	_, ok2 := svc.PutItem(put_input)
	if ok2 != nil {
		return ok2
	}
	get_adviser_input := &dynamodb.GetItemInput{
		TableName: aws.String("adviser_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(*get_output.Item["phone_number"].S)},
		},
	}
	adviser_info, ok3 := svc.GetItem(get_adviser_input)
	if ok3 != nil {
		return ok3
	}
	score, _ := strconv.ParseFloat(*adviser_info.Item["score"].N, 32)
	commentnum, _ := strconv.Atoi(*adviser_info.Item["commentnum"].N)
	if commentnum != 0 {
		commentnum++
	}
	newscore := (score*float64(commentnum) + float64(comment_request.Score)) / float64(commentnum)
	update_input := &dynamodb.UpdateItemInput{
		TableName: aws.String("adviser_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(comment.PhoneNumber)},
		},
		UpdateExpression: aws.String("SET score=:val1,commentnum=commentnum+:val2"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val1": {N: aws.String(fmt.Sprintf("%f", newscore))},
			":val2": {N: aws.String("1")},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}
	_, ok4 := svc.UpdateItem(update_input)
	if ok4 != nil {
		return ok4
	}
	return nil
}
