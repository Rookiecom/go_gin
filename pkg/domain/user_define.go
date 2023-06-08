package domain

type Loginreq struct {
	PhoneNumber string `json:"phone_number" dynamodbbav:"phone_number" `
	Password    string `json:"password" dynamodbbav:"password"`
}

type UserEdit struct {
	Name   string `json:"name" dynamodbbav:"name"`
	Birth  string `json:"birth" dynamodbbav:"birth"`
	Gender string `json:"gender" dynamodbbav:"gender"`
	Bio    string `json:"bio" dynamodbbav:"bio"`
	About  string `json:"about" dynamodbbav:"about"`
}
type UserDisEdit struct {
	Coin int    `json:"coin" dynamodbbav:"coin"`
	Uuid string `json:"uuid" dynamodbav:"uuid"`
}
type UserRequest struct {
	Loginreq
	UserEdit
}
type UserRegister struct {
	UserRequest
	UserDisEdit
}
type UserHomePage struct {
	UserEdit
	UserDisEdit
}
type UserRequestUpdate struct {
	IsEdit bool `json:"is_edit" dynamodbbav:"is_edit"`
	UserEdit
}
type AdviserCommentRequeat struct {
	AdviserId string `json:"adviser_id" dynamodbav:"adviser_id"`
	Content   string `json:"content" dynaomdodbav:"content"`
	Score     int    `json:"score" dynamodbav:"score"`
}
type AdviserComment struct {
	PhoneNumber string  `json:"phone_number" dynamodbav:"phone_number"`
	Content     string  `json:"content" dynaomdodbav:"content"`
	Score       float64 `json:"score" dynamodbav:"score"`
	Time        string  `json:"time" dynamodbav:"time"`
	Commentor   string  `json:"commentor" dynamodbav:"commentor"`
}
