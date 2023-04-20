package dingrobot

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

const (
	mt_markdown = "markdown"
)

type markdown struct {
	Title string `json:"title"` // required
	Text  string `json:"text"`  // required
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
	timestamp = time.Now().UnixNano() / 1e6
	str := fmt.Sprintf("%d\n%s", timestamp, p.signSecret)
	h := hmac.New(sha256.New, []byte(p.signSecret))
	h.Write([]byte(str))
	sign = base64.StdEncoding.EncodeToString(h.Sum(nil))
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
		"msgtype": mt_markdown,
		"markdown": &markdown{
			Title: title,
			Text:  msg,
		},
	}).Post(address.String())
	if err != nil {
		return
	}
	v, err := fastjson.Parse(resp.String())
	if err != nil {
		return
	}
	code, err := v.Get("errcode").Int64()
	if err != nil {
		return
	}
	if code != 0 {
		err = errors.New(v.Get("errmsg").String())
		return
	}
	return
}
