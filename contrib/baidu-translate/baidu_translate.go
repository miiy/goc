package baidu_translate

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type BaiduTranslate struct {
	appId  string
	secKey string
}

const (
	Auto = "auto"
	ZH   = "zh"
	EN   = "en"
)

var TransApiHost = "https://fanyi-api.baidu.com/api/trans/vip/translate"

type TransResponse struct {
	From        string        `json:"from"`
	To          string        `json:"to"`
	TransResult []TransResult `json:"trans_result"`
	ErrorCode   string        `json:"error_code"`
	ErrorMsg    string        `json:"error_msg"`
}

type TransResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func NewBaiduTranslate(appId, secKey string) *BaiduTranslate {
	return &BaiduTranslate{
		appId:  appId,
		secKey: secKey,
	}
}

// Translate
// http://api.fanyi.baidu.com/product/113
func (t *BaiduTranslate) Translate(q string, from, to string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	if len(from) == 0 {
		from = "auto"
	}
	if len(to) == 0 {
		return "", errors.New("invalid to")
	}

	min := 10000
	max := 99999
	salt := strconv.Itoa(rand.Intn(max-min) + min)
	sign := buildSign(t.appId, q, salt, t.secKey)

	reqData := url.Values{
		"appid": []string{t.appId},
		"q":     []string{q},
		"from":  []string{from},
		"to":    []string{to},
		"salt":  []string{salt},
		"sign":  []string{sign},
	}
	req, err := http.NewRequest(http.MethodPost, TransApiHost, strings.NewReader(reqData.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tr TransResponse
	err = json.Unmarshal(respBody, &tr)
	if err != nil {
		return "", err
	}
	fmt.Printf("%+v", tr)

	if tr.ErrorCode != "" {
		return "", errors.New(tr.ErrorMsg)
	}

	if len(tr.TransResult) > 0 {
		return tr.TransResult[0].Dst, nil
	}

	return "", nil
}

func buildSign(appId, q, salt, secKey string) string {
	h := md5.New()
	h.Write([]byte(appId + q + salt + secKey))
	return hex.EncodeToString(h.Sum(nil))
}
