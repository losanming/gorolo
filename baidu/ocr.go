package baidu

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Token struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

const (
	// 客户端凭证类型，固定为client_credentials
	grantType string = "client_credentials"
	// 应用的API Key
	apiKey string = "OsPafGjabn3gakGMtvKiKb5y"
	// 应用的Secret Key
	secretKey string = "KHcOQt9FhGDHtQ0L1ifmfFXMgMAzEicj"
	// 授权服务地址
	tokenUrl string = "https://aip.baidubce.com/oauth/2.0/token"
	// 文字识别（高精度）API接口地址
	accurateBasicUrl string = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate"
)

type accessToken struct {
	RefreshToken  string `json:"refresh_token"`
	ExpiresIn     uint32 `json:"expires_in"`
	Scope         string `json:"scope"`
	SessionKey    string `json:"session_key"`
	AccessToken   string `json:"access_token"`
	SessionSecret string `json:"session_secret"`
}

type WordsResult struct {
	Words    string    `json:"words"`
	Location Locations `json:"location"`
}

type Locations struct {
	Top    uint32 `json:"top"`
	Left   uint32 `json:"left"`
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
}

type Words struct {
	WordsResult    []WordsResult `json:"words_result"`
	LogId          uint64        `json:"log_id"`
	WordsResultNum uint32        `json:"words_result_num"`
}

// 获取access_token
// access_token有效期一般是30天
// 官方文档：https://ai.baidu.com/ai-doc/REFERENCE/Ck3dwjhhu
func getAccessToken() (data accessToken) {
	requestUrl := fmt.Sprintf("%s?grant_type=%s&client_id=%s&client_secret=%s", tokenUrl, grantType, apiKey, secretKey)
	response, e := http.Get(requestUrl)
	defer response.Body.Close()
	printError(e)
	s := unCoding(response)
	e = json.Unmarshal([]byte(s), &data)
	printError(e)
	return
}

// 文字识别（高精度）
// 官方文档：https://cloud.baidu.com/doc/OCR/s/1k3h7y3db
// img 图片地址
func AccurateBasic(img string) (data Words) {
	requestUrl := fmt.Sprintf("%s?access_token=%s", accurateBasicUrl, getAccessToken().AccessToken)
	f, e := os.Open(img)
	printError(e)
	defer f.Close()
	d, e := ioutil.ReadAll(f)
	printError(e)
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var coder = base64.NewEncoding(base64Table)
	imgString := coder.EncodeToString(d)
	printError(e)
	values := url.Values{"image": {imgString}}
	response, e := http.Post(requestUrl, "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
	defer response.Body.Close()
	s := unCoding(response)
	e = json.Unmarshal([]byte(s), &data)
	printError(e)
	return
}

func unCoding(r *http.Response) (body string) {
	if r.StatusCode == 200 {
		switch r.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ := gzip.NewReader(r.Body)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)
				if err != nil && err != io.EOF {
					panic(err)
				}
				if n == 0 {
					break
				}
				body += string(buf)
			}
		default:
			bodyByte, _ := ioutil.ReadAll(r.Body)
			body = string(bodyByte)
		}
	} else {
		bodyByte, _ := ioutil.ReadAll(r.Body)
		body = string(bodyByte)
	}
	return
}

func printError(e error) {
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}
}
