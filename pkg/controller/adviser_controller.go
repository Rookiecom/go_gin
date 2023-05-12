package controller

import (
	"fmt"
	"gowork/pkg/domain"
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
	c.String(http.StatusOK, "请输入以下信息完成顾问注册\n"+
		"phone_number password name  coin workstate ordernum score commentnum\n")
	var adviser_request domain.AdviserRequest
	c.ShouldBindJSON(&adviser_request)
	adviser_register := domain.AdviserRegister{
		AdviserRequest: adviser_request,
		AdviserDisEdit: *domain.AdviserDisEditInit(),
	}
	adviser_item, _ := dynamodbattribute.MarshalMap(adviser_register)
	adviser_item_input := &dynamodb.PutItemInput{
		Item:      adviser_item,
		TableName: aws.String("adviser"),
	}
	svc := domain.Svc
	_, err := svc.PutItem(adviser_item_input)
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
			TableName: aws.String("adviser"),
		}
		reqoutput, loginerr := svc.GetItem(reqinput)
		if loginerr == nil && reqoutput != nil {
			if *reqoutput.Item["password"].S == loginreq.Password {
				generateToken(c, "adviser", loginreq.PhoneNumber, 30)
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
		TableName: aws.String("adviser"),
		Key: map[string]*dynamodb.AttributeValue{
			"phone_number": {S: aws.String(claim.PhoneNumber)},
		},
	}
	reqoutput, _ := svc.GetItem(reqinput)
	fmt.Println(reqoutput)
	var home_page_inf domain.AdviserHomePage
	dynamodbattribute.UnmarshalMap(reqoutput.Item, &home_page_inf)
	fmt.Println(home_page_inf)
	c.JSON(http.StatusOK, home_page_inf)
}
func (actl *adviser_controller) AdviserInformationEdit(c *gin.Context) {
	ClaimsFormContext, _ := c.Get(util.GinContextKey)
	claim := ClaimsFormContext.(*myjwt.CustomClaims)
	svc := domain.Svc
	reqinput := &dynamodb.GetItemInput{
		TableName: aws.String("adviser"),
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
		fmt.Println(*update_req)
		if update_req.IsEdit {
			update_output, err := update_req.Update(claim.PhoneNumber)

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
