package domain

// Account 登录的12306的账户信息结构
type Account struct {
	Username string
	AppToken string
	Logined  bool
}

// Passenger 乘客信息
type Passenger struct {
	Code                string `json:"code"`
	PassengerName       string `json:"passenger_name"`
	SexCode             string `json:"sex_code"`
	BornDate            string `json:"born_date"`
	CountryCode         string `json:"country_code"`
	PassengerIDTypeCode string `json:"passenger_id_type_code"`
	PassengerIDTypeName string `json:"passenger_id_type_name"`
	PassengerIDNo       string `json:"passenger_id_no"`
	PassengerType       string `json:"passenger_type"`
	PassengerFlag       string `json:"passenger_flag"`
	PassengerTypeName   string `json:"passenger_type_name"`
	MobileNo            string `json:"mobile_no"`
	PhoneNo             string `json:"phone_no"`
	Email               string `json:"email"`
	Address             string `json:"address"`
	Postalcode          string `json:"postalcode"`
	FirstLetter         string `json:"first_letter"`
}
