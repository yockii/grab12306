package captcha

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const ruokuaiCaptchaURL = "http://api.ruokuai.com/create.json"
const ruokuaiUsername = "yockii"
const ruokuaiPassword = "1986327"
const ruokuaiTypeid = "6113"
const ruokuaiSoftid = "1"
const ruokuaiSoftkey = "b40ffbee5c1cf4e38028c197eb2fc751"
const ruokuaiTimeout = ""

var ruokuaiPwd = md5Encode([]byte(ruokuaiPassword))

// RuokuaiResult 若快返回的结果
type RuokuaiResult struct {
	Result string
	ID     string `json:"Id"`
}

// RuokuaiError 若快返回的错误信息
type RuokuaiError struct {
	Error     string
	ErrorCode string `json:"Error_Code"`
	Request   string
}

func md5Encode(v []byte) string {
	h := md5.New()
	h.Write(v)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// RuokuaiCaptchaResult 获取若快打码返回的结果，结果123表示第一排1、2、3三个图片
func RuokuaiCaptchaResult(imgContent []byte) (result string, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	bodyWriter.WriteField("username", ruokuaiUsername)
	bodyWriter.WriteField("password", ruokuaiPwd)
	bodyWriter.WriteField("typeid", ruokuaiTypeid)
	bodyWriter.WriteField("softid", ruokuaiSoftid)
	bodyWriter.WriteField("softkey", ruokuaiSoftkey)
	// bodyWriter.WriteField("timeout", timeout)

	fileWriter, err := bodyWriter.CreateFormFile("image", "filename")
	if err != nil {
		return
	}

	io.Copy(fileWriter, bytes.NewReader(imgContent))

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	res, err := http.Post(ruokuaiCaptchaURL, contentType, bodyBuf)
	if err != nil {
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var ruokuaiResult RuokuaiResult
	err = json.Unmarshal(content, &ruokuaiResult)
	if err != nil {
		return
	}

	if ruokuaiResult.Result == "" {
		var ruokuaiError RuokuaiError
		err = json.Unmarshal(content, &ruokuaiError)
		if err != nil {
			return
		}
		err = errors.New(ruokuaiError.Error)
		return
	}
	result = ruokuaiResult.Result
	return
}
