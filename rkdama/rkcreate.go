package rkdama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type GoResult struct {
	Result string
	Id     string
}

type ErrorCreate struct {
	Error      string
	Error_Code string
	Request    string
}

//JSON：{"Error":"错误提示信息","Error_Code":"","Request":""}

func RKCreate(username string, password string, typeid string, timeout string, softid string, softkey string, filename string) (*GoResult, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	bodyWriter.WriteField("username", username)
	bodyWriter.WriteField("password", password)
	bodyWriter.WriteField("typeid", typeid)
	bodyWriter.WriteField("timeout", timeout)
	bodyWriter.WriteField("softid", softid)
	bodyWriter.WriteField("softkey", softkey)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("image", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	defer fh.Close()

	//copy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}

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

	var goresult GoResult
	err1 := json.Unmarshal(resp_body, &goresult)
	if err1 != nil {
		return nil, err1
	}
	var errorcreate ErrorCreate
	if goresult.Result == "" {

		err2 := json.Unmarshal(resp_body, &errorcreate)
		if err2 != nil {
			return nil, err2
		}
		goresult.Result = errorcreate.Error
	}

	return &goresult, nil
}
