package domain

type AdviserEdit struct {
	Name      string `json:"name" dynamodbbav:"name"`
	WorkState bool   `json:"workstate" dynamodbbav:"workstate"`
}

type AdviserDisEdit struct {
	Coin           int     `json:"coin" dynamodbbav:"coin"`
	OrderNum       int     `json:"ordernum" dynamodbbav:"ordernum"`
	OrderFinishNum int     `json:"order_finish_num" dynamodbav:"order_finish_num"`
	Score          float64 `json:"score" dynamodbbav:"score"`
	CommentNum     int     `json:"commentnum" dynamodbbav:"commentnum"`
	Uuid           string  `json:"uuid" dynamodbav:"uuid"`
}
type AdviserRequest struct {
	Loginreq
	AdviserEdit
}
type AdviserRegister struct {
	AdviserRequest
	AdviserDisEdit
}
type AdviserHomePage struct {
	AdviserEdit
	AdviserDisEdit
}
type AdviserRequestUpdate struct {
	IsEdit bool `json:"is_edit" dynamodbbav:"is_edit"`
	AdviserEdit
}
