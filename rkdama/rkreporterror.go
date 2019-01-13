package rkdama

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type ReportResult struct {
	Result string
}

type ErrorReport struct {
	Error      string
	Error_Code string
	Request    string
}

func RKReportError(username string, userpassword string, softid string, softkey string, id string) (*ReportResult, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	bodyWriter.WriteField("username", username)
	bodyWriter.WriteField("password", userpassword)
	bodyWriter.WriteField("softid", softid)
	bodyWriter.WriteField("softkey", softkey)
	bodyWriter.WriteField("id", id)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post("http://api.ruokuai.com/reporterror.json", contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reportresult ReportResult
	err1 := json.Unmarshal(resp_body, &reportresult)
	if err1 != nil {
		return nil, err1
	}

	var errorreport ErrorReport
	if reportresult.Result == "" {
		err2 := json.Unmarshal(resp_body, &errorreport)
		if err2 != nil {
			return nil, err2
		}
		reportresult.Result = errorreport.Error
	}

	return &reportresult, nil
}
