package major

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"github.com/yockii/grab12306/config"
	netutil "github.com/yockii/grab12306/utils/net"
)

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

// MyPassengers 用户账号中的乘客信息
var MyPassengers []*Passenger

// FetchMyPassengers 从12306获取乘客信息
func FetchMyPassengers() (success bool) {
	pageIndex := 1
	pageSize := 10
	var myPs []*Passenger
	success, ps, totalPage := fetchPassengers(pageIndex, pageSize)
	if success {
		if len(ps) > 0 {
			for _, v := range ps {
				myPs = append(myPs, v)
			}
			if totalPage > 1 {
				for i := 1; i < totalPage; i++ {
					_, ps, _ = fetchPassengers(pageIndex+i, pageSize)
					for _, v := range ps {
						myPs = append(myPs, v)
					}
				}
			}
		}
		MyPassengers = myPs
	} else {
		fmt.Println("获取乘客信息失败")
	}
	return
}

func fetchPassengers(pageIndex, pageSize int) (success bool, ps []*Passenger, totalPage int) {
	data := make(url.Values)
	data["pageIndex"] = []string{strconv.Itoa(pageIndex)}
	data["pageSize"] = []string{strconv.Itoa(pageSize)}

	res, err := netutil.GetMajorClient().PostForm(config.Urls["passengers"], data)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	status := jsoniter.Get(content, "status").ToBool()
	if !status {
		fmt.Println("乘客信息获取失败", string(content[:]))
		return
	}
	resData := jsoniter.Get(content, "data")
	totalPage = resData.Get("pageTotal").ToInt()
	datasStr := resData.Get("datas").ToString()
	jsoniter.UnmarshalFromString(datasStr, &ps)
	success = true
	return
}
