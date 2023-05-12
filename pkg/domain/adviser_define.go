package domain

type AdviserEdit struct {
	Name      string `json:"name" dynamodbbav:"name"`
	WorkState bool   `json:"workstate" dynamodbbav:"workstate"`
}

type AdviserDisEdit struct {
	Coin       int `json:"coin" dynamodbbav:"coin"`
	OrderNum   int `json:"ordernum" dynamodbbav:"ordernum"`
	Score      int `json:"score" dynamodbbav:"score"`
	CommentNum int `json:"commentnum" dynamodbbav:"commentnum"`
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

func AdviserDisEditInit() *AdviserDisEdit {
	return &AdviserDisEdit{
		Coin:       0,
		OrderNum:   0,
		Score:      0,
		CommentNum: 0,
	}
}
