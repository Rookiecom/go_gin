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

type adviser_controller struct {
}

func NewAdviserController() *adviser_controller {
	return &adviser_controller{}
}
func (actl *adviser_controller) AdviserRegister(c *gin.Context) {
	klog.Infof("adviser register")
	ok1 := mydb.CreateAdviserTable()
	if ok1 != nil {
		c.String(http.StatusOK, ok1.Error()+"\n")
		return
	}
	ok2 := mydb.CreateIdTable()
	if ok2 != nil {
		c.String(http.StatusOK, ok2.Error()+"\n")
		return
	}
	c.String(http.StatusOK, "请输入以下信息完成顾问注册\n"+
		"phone_number password name  coin workstate ordernum score commentnum\n")
	var adviser_request domain.AdviserRequest
	c.ShouldBindJSON(&adviser_request)
	err := mydb.AdviserRegisterDB(&adviser_request)
	if err == nil {
		c.String(http.StatusOK, "注册成功")
	} else {
		c.String(http.StatusOK, "注册失败"+err.Error())
	}
}
func (actl *adviser_controller) AdviserLogin(c *gin.Context) {
	klog.Infof("adviser login for token")
	c.String(http.StatusOK, "填写如下信息登录：\nphone_number password name\n")
	var loginreq domain.Loginreq
	err := c.ShouldBindJSON(&loginreq)
	if err == nil {
		svc := domain.Svc
		reqinput := &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"phone_number": {S: aws.String(loginreq.PhoneNumber)},
			},
			TableName: aws.String("adviser_info"),
		}
		reqoutput, loginerr := svc.GetItem(reqinput)
		if loginerr == nil && reqoutput != nil {
			if *reqoutput.Item["password"].S == loginreq.Password {
				myjwt.GenerateToken(c, "adviser", loginreq.PhoneNumber, *reqoutput.Item["name"].S, 30)
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
func (actl *adviser_controller) AdviserHomePage(c *gin.Context) {
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claim := ClaimsFormContext.(*myjwt.CustomClaims)
	svc := domain.Svc
	reqinput := &dynamodb.GetItemInput{
		TableName: aws.String("adviser_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(claim.PhoneNumber)},
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
			":val": {S: aws.String(claim.PhoneNumber)},
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
func (actl *adviser_controller) AdviserInformationEdit(c *gin.Context) {
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claim := ClaimsFormContext.(*myjwt.CustomClaims)
	svc := domain.Svc
	reqinput := &dynamodb.GetItemInput{
		TableName: aws.String("adviser_info"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(claim.PhoneNumber)},
		},
	}
	reqoutput, _ := svc.GetItem(reqinput)
	var inf domain.AdviserEdit
	dynamodbattribute.UnmarshalMap(reqoutput.Item, &inf)
	c.String(http.StatusOK, "修改以前的信息\n")
	c.JSON(http.StatusOK, inf)
	var update_req = &domain.AdviserRequestUpdate{
		IsEdit:      false,
		AdviserEdit: inf,
	}
	if ok := c.ShouldBindJSON(update_req); ok == nil {
		if update_req.IsEdit {
			update_output, err := mydb.AdviserUpdate(update_req, claim.PhoneNumber)
			if err == nil {
				var new_inf domain.AdviserEdit
				dynamodbattribute.UnmarshalMap(update_output.Attributes, &new_inf)
				c.String(http.StatusOK, "\n修改以后的信息\n")
				c.JSON(http.StatusOK, new_inf)
			} else {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("edit 传送失败")
		}

	} else {
		c.String(http.StatusOK, ok.Error())
	}
}
func (actl *adviser_controller) TakeOrder(c *gin.Context) {
	klog.Info("adviser token order")
	orderid := c.Query("id")
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claim := ClaimsFormContext.(*myjwt.CustomClaims)
	err := mydb.OrderTake(orderid, claim.Name, claim.PhoneNumber)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	} else {
		c.String(http.StatusOK, "succeeded to take order\n")
	}
}
func (actl *adviser_controller) AdviserGetFlow(c *gin.Context) {
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claims := ClaimsFormContext.(*myjwt.CustomClaims)
	page_request := &domain.Pagination{}
	c.ShouldBindJSON(page_request)
	query_output, err := mydb.GetFlow(claims.PhoneNumber, page_request)
	if err != nil {
		c.String(http.StatusOK, "fail to get coin flow")
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
