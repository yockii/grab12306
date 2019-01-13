package rkdama

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type RegisterResult struct {
	Result string
}

type ErrorRegister struct {
	Error      string
	Error_Code string
	Request    string
}

func RKRegister(username string, userpassword string, email string) (*RegisterResult, error) {

	senddata := "username=$1&password=$2&email=$3"
	senddata = strings.Replace(senddata, "$1", username, -1)
	senddata = strings.Replace(senddata, "$2", userpassword, -1)
	senddata = strings.Replace(senddata, "$3", email, -1)
	postdata := strings.NewReader(senddata)
	contentType := "application/x-www-form-urlencoded"
	resp, err := http.Post("http://api.ruokuai.com/register.json", contentType, postdata)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var registerresult RegisterResult
	err1 := json.Unmarshal(resp_body, &registerresult)
	if err1 != nil {
		return nil, err1
	}

	var errorregister ErrorRegister
	if registerresult.Result == "" {
		err2 := json.Unmarshal(resp_body, &errorregister)
		if err2 != nil {
			return nil, err2
		}
		registerresult.Result = errorregister.Error
	}

	return &registerresult, nil
}
