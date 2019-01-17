package domain

// UamauthClientResponse 获取校验后的用户名返回信息
type UamauthClientResponse struct {
	ResultCode int    `json:"result_code"`
	Username   string `json:"username"`
}

// AuthResponse 校验请求返回信息，主要获取token
type AuthResponse struct {
	ResultCode int    `json:"result_code"`
	AppToken   string `json:"newapptk"`
}

// BaseLoginResponse 登录返回
type BaseLoginResponse struct {
	ResultCode int    `json:"result_code"`
	Uamtk      string `json:"uamtk"`
}

// CheckCaptchaResponse 登录返回
type CheckCaptchaResponse struct {
	ResultCode    int    `json:"result_code"`
	ResultMessage string `json:"result_message"`
}

// SubmitOrderResponse 提交订单信息返回
type SubmitOrderResponse struct {
	ValidateMessagesShowID string `json:"validateMessagesShowId"`
	Status                 bool   `json:"status"`
	HTTPStatus             int    `json:"httpstatus"`
	Data                   string `json:"data"`
	Messages               string `json:"messages"`
}

// CheckOrderInfoResponse checkOrderInfo的返回信息
type CheckOrderInfoResponse struct {
	ValidateMessagesShowID string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		IfShowPassCode     string `json:"ifShowPassCode"`
		CanChooseBeds      string `json:"canChooseBeds"`  // 是否可以选铺位
		CanChooseSeats     string `json:"canChooseSeats"` // 是否可以选座位
		ChooseSeats        string `json:"choose_Seats"`   // 可选座位类型
		IsCanChooseMid     string `json:"isCanChooseMid"`
		IfShowPassCodeTime string `json:"ifShowPassCodeTime"`
		SubmitStatus       bool   `json:"submitStatus"` // 订单检测结果
		SmokeStr           string `json:"smokeStr"`
		ErrMsg             string `json:"errMsg"` // submitStatus=false时的错误信息
	} `json:"data"`
}

// AutoSubmitOrderResponse autoSubmitOrderRequest的返回信息
type AutoSubmitOrderResponse struct {
	ValidateMessagesShowID string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		IfShowPassCode     string `json:"ifShowPassCode"`
		CanChooseBeds      string `json:"canChooseBeds"`  // 是否可以选铺位
		CanChooseSeats     string `json:"canChooseSeats"` // 是否可以选座位
		ChooseSeats        string `json:"choose_Seats"`   // 可选座位类型
		IsCanChooseMid     string `json:"isCanChooseMid"`
		IfShowPassCodeTime string `json:"ifShowPassCodeTime"`
		SubmitStatus       bool   `json:"submitStatus"` // 订单检测结果
		SmokeStr           string `json:"smokeStr"`
		Result             string `json:"result"`
	} `json:"data"`
}

// LeftTicketResponse 余票查询返回信息
type LeftTicketResponse struct {
	HTTPStatus int    `json:"httpstatus"`
	Messages   string `json:"messages"`
	Status     bool   `json:"status"`
	Data       struct {
		Flag   string            `json:"flag"`
		Map    map[string]string `json:"map"`
		Result []string          `json:"result"`
	} `json:"data"`
}

// GetQueueCountResponse 获取余票及队列信息返回
type GetQueueCountResponse struct {
	ValidateMessagesShowID string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		Count  string `json:"count"`  // 排队人数
		Ticket string `json:"ticket"` // 余票数
		Op2    string `json:"op_2"`
		CountT string `json:"countT"`
		Op1    string `json:"op_1"`
	} `json:"data"`
}

// ConfirmQueueResponse confirmQueue返回
type ConfirmQueueResponse struct {
	ValidateMessagesShowID string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		SubmitStatus bool `json:"submitStatus"`
	} `json:"data"`
}

// QueryOrderResponse 查询订单返回
type QueryOrderResponse struct {
	ValidateMessagesShowID string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	HTTPStatus             int      `json:"httpstatus"`
	Messages               []string `json:"messages"`
	Data                   struct {
		QueryOrderWaitTimeStatus string  `json:"queryOrderWaitTimeStatus"`
		Count                    int     `json:"count"`
		WaitTime                 int     `json:"waitTime"`
		RequestID                int64   `json:"requestId"`
		WaitCount                int     `json:"waitCount"`
		TourFlag                 string  `json:"tourFlag"`
		OrderID                  Nstring `json:"orderId"`
	} `json:"data"`
}
