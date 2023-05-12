package controller

import (
	"fmt"
	domain "gowork/pkg/domain"
	myjwt "gowork/pkg/middleware"
	"gowork/pkg/util"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/dgrijalva/jwt-go"
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

	var userreq domain.UserRequest
	c.String(200, "请填写以下信息完成用户注册\n")
	c.String(200, "phone_number password name birth gender bio about\n")
	user_register := domain.UserRegister{
		UserRequest: userreq,
		UserDisEdit: *domain.UserDisEditInit(),
	}
	c.ShouldBindJSON(&user_register)
	useritem, _ := dynamodbattribute.MarshalMap(user_register)
	useriteminput := &dynamodb.PutItemInput{
		Item:      useritem,
		TableName: aws.String("user"),
	}
	svc := domain.Svc
	_, err := svc.PutItem(useriteminput)
	if err == nil {
		c.String(http.StatusOK, "注册成功")
	} else {
		c.String(http.StatusOK, "注册失败"+err.Error())
	}
}
func (uctl *user_controller) UserLogin(c *gin.Context) {
	klog.Infof("user login for token")
	c.String(http.StatusOK, "填写如下信息登录：\nphone_number password name\n")
	var loginreq domain.Loginreq
	err := c.ShouldBindJSON(&loginreq)
	if err == nil {
		svc := domain.Svc
		reqinput := &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"phone_number": {S: aws.String(loginreq.PhoneNumber)},
			},
			TableName: aws.String("user"),
		}
		reqoutput, loginerr := svc.GetItem(reqinput)
		if loginerr == nil && reqoutput != nil {
			if *reqoutput.Item["password"].S == loginreq.Password {
				generateToken(c, "user", loginreq.PhoneNumber, 30)
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
	svc := domain.Svc
	reqinput := &dynamodb.GetItemInput{
		TableName: aws.String("user"),
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
		TableName: aws.String("user"),
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
			update_output, err := update_req.Update(claim.PhoneNumber)

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
func generateToken(c *gin.Context, roleid string, phone_number string, ExpireTimeByMinute int) {
	j := myjwt.NewJwt()

	claims := myjwt.CustomClaims{
		Role:        roleid,
		PhoneNumber: phone_number,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),
			ExpiresAt: int64(time.Now().Unix() + int64(ExpireTimeByMinute*60)),
			Issuer:    "rookie",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "登录成功",
		"token":  token,
	})
}
