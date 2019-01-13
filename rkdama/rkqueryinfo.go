package rkdama

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type InfoResult struct {
	Score        string
	HistoryScore string
	TotalScore   string
	TotalTopic   string
}

type ErrorInfo struct {
	Error      string
	Error_Code string
	Request    string
}

func RKQueryInfo(username string, userpassword string) (*InfoResult, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	bodyWriter.WriteField("username", username)
	bodyWriter.WriteField("password", userpassword)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post("http://api.ruokuai.com/info.json", contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var inforesult InfoResult
	err1 := json.Unmarshal(resp_body, &inforesult)
	if err1 != nil {
		return nil, err1
	}

	var errorinfo ErrorInfo
	if inforesult.Score == "" {
		err2 := json.Unmarshal(resp_body, &errorinfo)
		if err2 != nil {
			return nil, err2
		}
		inforesult.Score = errorinfo.Error
	}

	return &inforesult, err
}
