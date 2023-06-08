package controller

import (
	"fmt"
	mydb "gowork/pkg/db"
	domain "gowork/pkg/domain"
	myjwt "gowork/pkg/middleware"
	"gowork/pkg/util"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

type user_controller struct {
	//user sevice
}

func NewUserController() *user_controller {
	return &user_controller{}
}

func (uctl *user_controller) UserRegister(c *gin.Context) {
	klog.Infof("user register")
	ok1 := mydb.CreateUserTable()
	if ok1 != nil {
		c.String(http.StatusOK, ok1.Error()+"\n")
		return
	}
	ok2 := mydb.CreateIdTable()
	if ok2 != nil {
		c.String(http.StatusOK, ok2.Error()+"\n")
		return
	}
	c.String(200, "请填写以下信息完成用户注册\n")
	c.String(200, "phone_number password name birth gender bio about\n")
	var userreq domain.UserRequest
	dataok := c.ShouldBindJSON(&userreq)
	if dataok != nil {
		c.String(http.StatusOK, "网页传输数据失败")
		return
	}
	err := mydb.UserRegisterDB(&userreq)
	if err == nil {
		c.String(http.StatusOK, "注册成功")
	} else {
		c.String(http.StatusOK, "注册失败"+err.Error())
	}
}
func (uctl *user_controller) UserLogin(c *gin.Context) {

	klog.Infof("user login for token")
	_, asd := c.Get(util.GinContextKey)
	if !asd {
		fmt.Println("fafadfafdfadssdffafadsfasdfasdfa")
	}
	c.String(http.StatusOK, "填写如下信息登录：\nphone_number password name\n")
	var loginreq domain.Loginreq
	err := c.ShouldBindJSON(&loginreq)
	if err == nil {
		svc := domain.Svc
		reqinput := &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"phone_number": {S: aws.String(loginreq.PhoneNumber)},
			},
			TableName: aws.String("user_info"),
		}
		reqoutput, loginerr := svc.GetItem(reqinput)
		if loginerr == nil && reqoutput != nil {
			if *reqoutput.Item["password"].S == loginreq.Password {
				myjwt.GenerateToken(c, "user", loginreq.PhoneNumber, *reqoutput.Item["name"].S, 30)
			} else {
				fmt.Println(*reqoutput.Item["password"].S)
				fmt.Println(loginreq.Password)
				c.String(http.StatusOK, "登录信息错误")
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "数据库传递的数据解析失败",
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "网页传递的数据解析失败",
		})
	}
}
func (uctl *user_controller) UserHomePage(c *gin.Context) {
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claim := ClaimsFormContext.(*myjwt.CustomClaims)
	fmt.Println(claim)
	svc := domain.Svc
	reqinput := &dynamodb.GetItemInput{
		TableName: aws.String("user_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(claim.PhoneNumber)},
		},
	}
	reqoutput, _ := svc.GetItem(reqinput)
	var home_page_inf domain.UserHomePage
	dynamodbattribute.UnmarshalMap(reqoutput.Item, &home_page_inf)
	c.JSON(http.StatusOK, home_page_inf)
}
func (uctl *user_controller) UserInformationEdit(c *gin.Context) {
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claim := ClaimsFormContext.(*myjwt.CustomClaims)
	svc := domain.Svc
	reqinput := &dynamodb.GetItemInput{
		TableName: aws.String("user_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(claim.PhoneNumber)},
		},
	}
	reqoutput, _ := svc.GetItem(reqinput)
	var inf domain.UserEdit
	dynamodbattribute.UnmarshalMap(reqoutput.Item, &inf)
	c.String(http.StatusOK, "修改以前的信息\n")
	c.JSON(http.StatusOK, inf)
	var update_req = &domain.UserRequestUpdate{
		IsEdit:   false,
		UserEdit: inf,
	}
	if ok := c.ShouldBindJSON(update_req); ok == nil {
		fmt.Println(*update_req)
		if update_req.IsEdit {
			update_output, err := mydb.UserUpdate(update_req, claim.PhoneNumber)

			if err == nil {
				var new_inf domain.UserEdit
				dynamodbattribute.UnmarshalMap(update_output.Attributes, &new_inf)
				c.String(http.StatusOK, "\n修改以后的信息\n")
				c.JSON(http.StatusOK, new_inf)
			} else {
				fmt.Println(err.Error())
			}
		}
	} else {
		c.String(http.StatusOK, ok.Error())
	}
}
func (uctl *user_controller) GetAdviserList(c *gin.Context) {
	input := &dynamodb.ScanInput{
		TableName:            aws.String("adviser_info"),
		ProjectionExpression: aws.String("#phone, #name, #uuid"),
		ExpressionAttributeNames: map[string]*string{
			"#phone": aws.String("phone_number"),
			"#name":  aws.String("name"),
			"#uuid":  aws.String("uuid"),
		},
	}
	svc := domain.Svc
	result, _ := svc.Scan(input)
	for _, item := range result.Items {
		c.String(http.StatusOK, "name:"+*item["name"].S+"  uuid:"+*item["uuid"].S+"\n")
	}

}
func (uctl *user_controller) VisitAdviser(c *gin.Context) {
	adviserid := c.Query("uuid")

	phonereq := &dynamodb.GetItemInput{
		TableName: aws.String("all_id"),
		Key: map[string]*dynamodb.AttributeValue{
			"uuid": {S: aws.String(adviserid)},
		},
	}
	svc := domain.Svc
	result, _ := svc.GetItem(phonereq)
	reqinput := &dynamodb.GetItemInput{
		TableName: aws.String("adviser_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(*result.Item["phone_number"].S)},
		},
	}
	reqoutput, _ := svc.GetItem(reqinput)
	var home_page_inf domain.AdviserHomePage
	dynamodbattribute.UnmarshalMap(reqoutput.Item, &home_page_inf)
	c.JSON(http.StatusOK, home_page_inf)
	scan_input := &dynamodb.ScanInput{
		TableName:        aws.String("adviser_comment"),
		FilterExpression: aws.String("phone_number=:val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {S: aws.String(*result.Item["phone_number"].S)},
		},
		ProjectionExpression: aws.String("commentor,#t,score,content"),
		ExpressionAttributeNames: map[string]*string{
			"#t": aws.String("time"),
		},
	}
	scan_output, err := svc.Scan(scan_input)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		var comments []domain.CommentDisplay
		for i := 0; i < int(*scan_output.Count); i++ {
			var temp domain.CommentDisplay
			dynamodbattribute.UnmarshalMap(scan_output.Items[i], &temp)
			comments = append(comments, temp)
		}
		c.JSON(http.StatusOK, comments)
	}

}
func (uctl *user_controller) HostOrder(c *gin.Context) {
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claim := ClaimsFormContext.(*myjwt.CustomClaims)
	ok1 := mydb.CreateOrderTable()
	if ok1 != nil {
		c.String(http.StatusOK, "failed to create order table\n"+ok1.Error()+"\n")
		return
	}
	ok2 := mydb.CreateSelectOrder()
	if ok2 != nil {
		c.String(http.StatusOK, "failed to create select order table"+ok2.Error()+"\n")
	}
	order := &domain.Order{}
	order.OrderHost = claim.Name

	c.ShouldBindJSON(order)
	err := mydb.OrderRegister(order, claim.PhoneNumber)
	if err != nil {
		c.String(http.StatusOK, "failed to register order\n"+err.Error())
	} else {
		c.String(http.StatusOK, "succeeded to register order\n")
	}
}
func (uctl *user_controller) FinishOrder(c *gin.Context) {
	ok := mydb.CreateFlowTable()
	if ok != nil {
		c.String(http.StatusOK, "failed to create order table\n"+ok.Error()+"\n")
		return
	}
	var data domain.FinishRequest
	c.ShouldBindJSON(&data)
	if !data.IsFinish {
		return
	}
	err := mydb.OrderFinish(data.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		c.String(http.StatusOK, "succeeded to finish order")
	}
}
func (uctl *user_controller) CommentOrder(c *gin.Context) {
	comment := domain.CommentRequest{}
	c.ShouldBindJSON(&comment)
	err := mydb.OrderComment(&comment)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		c.String(http.StatusOK, "succeeded to comment order")
	}
}

func (uctl *user_controller) RewardOrder(c *gin.Context) {
	reward := domain.RewardRequest{}
	c.ShouldBindJSON(&reward)
	err := mydb.OrderReward(&reward)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		c.String(http.StatusOK, "succeeded to reward order")
	}
}
func (uctl *user_controller) CollectAdviser(c *gin.Context) {

	ok := mydb.CreateCollect()
	if ok != nil {
		c.String(http.StatusOK, "failed to create collect table\n"+ok.Error()+"\n")
		return
	}
	adviser_id := c.Query("uuid")
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claim := ClaimsFormContext.(*myjwt.CustomClaims)
	err := mydb.CollectAdviser(claim.PhoneNumber, adviser_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		c.String(http.StatusOK, "succeeded to collect "+adviser_id)
	}
}
func (uctl *user_controller) GetCollectList(c *gin.Context) {
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claim := ClaimsFormContext.(*myjwt.CustomClaims)
	svc := domain.Svc
	input := &dynamodb.GetItemInput{
		TableName: aws.String("user_collect_adviser"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(claim.PhoneNumber)},
		},
	}
	output, err := svc.GetItem(input)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		var advisrelist []string
		dynamodbattribute.Unmarshal(output.Item["adviser_list"], &advisrelist)
		c.JSON(http.StatusOK, advisrelist)
	}
}

func (uctl *user_controller) UserGetFlow(c *gin.Context) {
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claims := ClaimsFormContext.(*myjwt.CustomClaims)
	page_request := &domain.Pagination{}
	c.ShouldBindJSON(page_request)
	query_output, err := mydb.GetFlow(claims.PhoneNumber, page_request)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		var output []domain.FlowOutput
		for i := 0; i < int(*query_output.Count); i++ {
			var flow domain.FlowOutput
			dynamodbattribute.UnmarshalMap(query_output.Items[i], &flow)
			output = append(output, flow)
		}
		c.JSON(http.StatusOK, output)
	}
}

func (uctl *user_controller) CommentAdviser(c *gin.Context) {
	ok := mydb.CreateAdviserCommentTable()
	if ok != nil {
		c.String(http.StatusOK, "failed to create adviser_comment table\n"+ok.Error()+"\n")
		return
	}
	comment_request := &domain.AdviserCommentRequeat{}
	c.ShouldBindJSON(comment_request)
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claims := ClaimsFormContext.(*myjwt.CustomClaims)
	err := mydb.CommentAdviser(claims.PhoneNumber, comment_request)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": err.Error(),
		})
	} else {
		c.String(http.StatusOK, "succeeded to comment "+comment_request.AdviserId)
	}
}
