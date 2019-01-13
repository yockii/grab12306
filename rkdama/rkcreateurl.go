package rkdama

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type GoUrlResult struct {
	Result string
	Id     string
}

type ErrorCreateUrl struct {
	Error      string
	Error_Code string
	Request    string
}

//JSON：{"Error":"错误提示信息","Error_Code":"","Request":""}

func RKCreateUrl(username string, password string, typeid string, timeout string, softid string, softkey string, url string) (*GoUrlResult, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	bodyWriter.WriteField("username", username)
	bodyWriter.WriteField("password", password)
	bodyWriter.WriteField("typeid", typeid)
	bodyWriter.WriteField("timeout", timeout)
	bodyWriter.WriteField("softid", softid)
	bodyWriter.WriteField("softkey", softkey)
	bodyWriter.WriteField("imageurl", url)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post("http://api.ruokuai.com/create.json", contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gourlresult GoUrlResult
	err1 := json.Unmarshal(resp_body, &gourlresult)
	if err1 != nil {
		return nil, err1
	}
	var errorcreateurl ErrorCreateUrl
	if gourlresult.Result == "" {

		err2 := json.Unmarshal(resp_body, &errorcreateurl)
		if err2 != nil {
			return nil, err2
		}
		gourlresult.Result = errorcreateurl.Error
	}

	return &gourlresult, nil
}
