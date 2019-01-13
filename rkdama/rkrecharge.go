package rkdama

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type RechargeResult struct {
	Result string
}

type ErrorRecharge struct {
	Error      string
	Error_Code string
	Request    string
}

func RKRecharge(username string, id string, password string) (*RechargeResult, error) {

	senddata := "username=$1&id=$2&password=$3"
	senddata = strings.Replace(senddata, "$1", username, -1)
	senddata = strings.Replace(senddata, "$2", id, -1)
	senddata = strings.Replace(senddata, "$3", password, -1)
	postdata := strings.NewReader(senddata)
	contentType := "application/x-www-form-urlencoded"
	resp, err := http.Post("http://api.ruokuai.com/recharge.json", contentType, postdata)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rechargeresult RechargeResult
	err1 := json.Unmarshal(resp_body, &rechargeresult)
	if err1 != nil {
		return nil, err1
	}

	var errorrecharge ErrorRecharge
	if rechargeresult.Result == "" {
		err2 := json.Unmarshal(resp_body, &errorrecharge)
		if err2 != nil {
			return nil, err2
		}
		rechargeresult.Result = errorrecharge.Error
	}

	return &rechargeresult, nil
}
