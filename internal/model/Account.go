package model

// Account 登录的12306的账户信息结构
type Account struct {
	Id        string
	UserId    string
	Username  string
	CookieJar []byte
}

// Passenger 乘客信息
type Passenger struct {
	Id            string
	AccountId     string
	PassengerUuid string `json:"passenger_uuid"` // 12306的uuid
	AllEncStr     string `json:"allEncStr"`
	PassengerName string `json:"passenger_name"`
	PassengerType string `json:"passenger_type"` // 1-成人
}
