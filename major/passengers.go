package major

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"

	"github.com/yockii/grab12306/constant"
	"github.com/yockii/grab12306/domain"

	netutil "github.com/yockii/grab12306/utils/net"
)

// MyPassengers 用户账号中的乘客信息
var MyPassengers []domain.Passenger

// FetchMyPassengers 从12306获取乘客信息
func FetchMyPassengers() (success bool) {
	pageIndex := 1
	pageSize := 10
	var myPs []domain.Passenger
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

func fetchPassengers(pageIndex, pageSize int) (success bool, ps []domain.Passenger, totalPage int) {
	data := make(url.Values)
	data["pageIndex"] = []string{strconv.Itoa(pageIndex)}
	data["pageSize"] = []string{strconv.Itoa(pageSize)}

	res, err := netutil.GetMajorClient().PostForm(constant.Urls["passengers"], data)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	var pRes PassengerResponse
	err = json.Unmarshal(content, &pRes)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !pRes.Status {
		fmt.Println("乘客信息获取失败", string(content[:]))
		return
	}
	totalPage = pRes.Data.PageTotal
	ps = pRes.Data.Datas
	success = true
	return
}

// PassengerResponse 乘客信息的返回结构
type PassengerResponse struct {
	HTTPStatus             int                `json:"httpstatus"`
	Messages               []string           `json:"messages"`
	Status                 bool               `json:"status"`
	ValidateMessages       interface{}        `json:"validateMessages"`
	ValidateMessagesShowID string             `json:"validateMessagesShowId"`
	Data                   PassengerPagedInfo `json:"data"`
}

// PassengerPagedInfo 乘客信息返回结果的分页信息
type PassengerPagedInfo struct {
	Flag      bool               `json:"flag"`
	PageTotal int                `json:"pageTotal"`
	Datas     []domain.Passenger `json:"datas"`
}
