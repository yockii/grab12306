package captcha

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"
)

const ruokuaiCaptchaURL = "http://api.ruokuai.com/create.json"
const ruokuaiUsername = "yockii"
const ruokuaiPassword = "1986327"
const ruokuaiTypeid = "6113"
const ruokuaiSoftid = "1"
const ruokuaiSoftkey = "b40ffbee5c1cf4e38028c197eb2fc751"
const ruokuaiTimeout = ""

var ruokuaiPwd = md5Encode([]byte(ruokuaiPassword))

func md5Encode(v []byte) string {
	h := md5.New()
	h.Write(v)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func RuokuaiCaptchaResult(imgContent []byte) (result string) {
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
		fmt.Println(err)
		return
	}

	io.Copy(fileWriter, bytes.NewReader(imgContent))

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	res, err := http.Post(ruokuaiCaptchaURL, contentType, bodyBuf)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	result = jsoniter.Get(content, "Result").ToString()
	// id := jsoniter.Get(content, "Id").ToString()
	return
}

func ruokuaiCaptchaFileResult(filename string) (result string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	// bodyWriter.SetBoundary("-------------RK")

	bodyWriter.WriteField("username", ruokuaiUsername)
	bodyWriter.WriteField("password", ruokuaiPwd)
	bodyWriter.WriteField("typeid", ruokuaiTypeid)
	bodyWriter.WriteField("softid", ruokuaiSoftid)
	bodyWriter.WriteField("softkey", ruokuaiSoftkey)
	// bodyWriter.WriteField("timeout", timeout)

	fileWriter, err := bodyWriter.CreateFormFile("image", filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fileWriter.Write(imgContent)
	// io.Copy(fileWriter, bytes.NewReader(imgContent))
	fh, _ := os.Open(filename)
	defer fh.Close()
	io.Copy(fileWriter, fh)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	res, err := http.Post(ruokuaiCaptchaURL, contentType, bodyBuf)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	result = jsoniter.Get(content, "Result").ToString()
	// id := jsoniter.Get(content, "Id").ToString()
	fmt.Printf("识别结果: %s", string(content[:]))
	return
}
