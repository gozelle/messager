package feishurobot

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/gozelle/fastjson"
	"github.com/gozelle/resty"
)

type content struct {
	Tag  string `json:"tag"`
	Text string `json:"text"` // required
}

func NewRobot(webhook, signSecret string) *Robot {
	robot := &Robot{
		webhook:    webhook,
		signSecret: signSecret,
		http:       resty.New(),
	}
	return robot
}

type Robot struct {
	webhook        string
	signSecret     string
	http           *resty.Client
	titleFormatter func(messages []interface{}) string
}

func (p *Robot) SetTitleFormatter(format func(messages []interface{}) string) {
	p.titleFormatter = format
}

func (p *Robot) sign() (timestamp int64, sign string) {
	timestamp = time.Now().Unix()
	sign, _ = GenSign(p.signSecret, timestamp)
	return
}

func (p *Robot) Push(ctx context.Context, title, msg string) (err error) {
	timestamp, sign := p.sign()

	address, err := url.Parse(p.webhook)
	if err != nil {
		log.Println(err)
		return
	}
	query := address.Query()
	query.Add("timestamp", strconv.FormatInt(timestamp, 10))
	query.Add("sign", sign)
	address.RawQuery = query.Encode()
	resp, err := p.http.R().SetBody(map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title": title,
					"content": [][]content{
						{
							{
								Tag:  "text",
								Text: msg,
							},
						},
					},
				},
			},
		},
	}).Post(address.String())
	if err != nil {
		return
	}
	v, err := fastjson.Parse(resp.String())
	if err != nil {
		return
	}
	
	code, err := v.Get("code").Int64()
	if err != nil {
		return
	}
	if code != 0 {
		err = errors.New(v.Get("msg").String())
		return
	}
	return
}

func GenSign(secret string, timestamp int64) (string, error) {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
