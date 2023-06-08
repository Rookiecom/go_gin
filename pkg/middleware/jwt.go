package middleware

import (
	"errors"
	"fmt"
	"gowork/pkg/util"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

type JWT struct {
	SigningKey []byte
}
type CustomClaims struct {
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	jwt.StandardClaims
}

func NewJwt() *JWT {
	return &JWT{
		[]byte(util.SignKey),
	}
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}
func (j *JWT) ParseToken(tokenstring string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("格式不正确")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("过期")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("尚未生效")
			} else {
				return nil, errors.New("token非法")
			}
		}
	}
	claims, ok := token.Claims.(*CustomClaims)
	if ok {
		return claims, nil
	} else {
		return nil, errors.New("token非法")
	}
}

func JWTAuth(model int) gin.HandlerFunc {
	//model代表3种模式,为1代表user，为2代表adviser，为3代表普通模式
	//user模式需要携带token解析出来的role信息为user，同理adviser
	return func(c *gin.Context) {
		klog.Infof("reqURI:%v", c.Request.RequestURI)
		if (strings.Contains(c.Request.RequestURI, "login")) || (strings.Contains(c.Request.RequestURI, "register")) {
			return
		} else if c.Request.RequestURI == "/" {
			return
		}
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "未携带token",
			})
			c.Abort()
			return
		}
		klog.Infof("token:%s", token)
		j := NewJwt()
		claims, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			c.Abort()
			return
		}
		switch model {
		case 1:
			{
				if claims.Role != "user" {
					c.String(http.StatusOK, "user身份不符")
					fmt.Println(model)
					c.Abort()
					return
				}
			}
		case 2:
			{
				if claims.Role != "adviser" {
					c.String(http.StatusOK, "adviser身份不符")
					c.Abort()
					return
				}
			}
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set(util.GinContextKey, claims)
		//可以通过 c.Get(util.Gin_Context_Key) 获取到之前存储的值
	}
}

func GenerateToken(c *gin.Context, roleid string, phone_number string, name string, ExpireTimeByMinute int) {
	j := NewJwt()

	claims := CustomClaims{
		Role:        roleid,
		PhoneNumber: phone_number,
		Name:        name,
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
