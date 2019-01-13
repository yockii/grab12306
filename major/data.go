package major

// Account 登录的12306的账户信息结构
type Account struct {
	Username string
	AppToken string
	Logined  bool
}

// MyAccount 我的12306账户信息
var MyAccount Account
