package domain

type OrderDisplay struct {
	ID          string `json:"id" dynamodbav:"id"`
	Name        string `json:"name" dynamodbav:"name"`
	Price       int    `json:"price" dynamodbav:"price"`
	Description string `json:"description" dynamodbav:"description"`
	Status      string `json:"status" dynamodbav:"status"`
	OrderHost   string `json:"order_host" dynamodbav:"order_Host"`
	OrderTaker  string `json:"order_taker" dynamodbav:"order_taker"`
	CommentText string `json:"comment" dynamodbav:"comment"`
}
type Order struct {
	OrderDisplay
	OrderHostPhone  string `json:"order_host_phone" dynamodbav:"order_host_phone"`
	OrderTakerPhone string `json:"order_taker_phone" dynamodbav:"order_taker_phone"`
}
type AuxiliaryQuery struct {
	PhoneNumber string `json:"phone_number" dynamodbav:"phone_number"`
	Name        string `json:"name" dynamodbav:"name"`
	ID          string `json:"id" dynamodbav:"id"`
	Type        string `json:"type" dynamodbav:"type"`
}

type FinishRequest struct {
	IsFinish bool   `json:"is_finish"`
	ID       string `json:"id"`
}
type CommentRequest struct {
	CommentText string `json:"comment_text" `
	ID          string `json:"id"`
}
type RewardRequest struct {
	Reward int    `json:"reward"`
	ID     string `json:"id"`
}
