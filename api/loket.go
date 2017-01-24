package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/maps90/librarian"
	"github.com/maps90/librarian/cache"
	gr "github.com/parnurzeal/gorequest"
)

var conf map[string]string
var baseUrl string

const CACHE_DURATION = 5 * time.Minute

type Loket struct {
	BaseUrl      string
	UserName     string
	Password     string
	ApiKey       string
	Token        string
	Response     LoketResponse
	Body         string
	Error        error
	TokenExpired bool
	OnCache      bool
}

type LoketResponse struct {
	Code    uint16      `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}

func getConfig(key string) string {
	if _, ok := conf[key]; ok {
		return conf[key]
	}
	return ""
}

func (l *Loket) getAuth() *Loket {
	if !l.TokenExpired && len(l.Token) != 0 {
		return l
	}
	if len(l.UserName) == 0 || len(l.Password) == 0 || len(l.ApiKey) == 0 {
		return l
	}
	var errs []error

	body := fmt.Sprintf(`{"username": "%s","password": "%s","APIKEY": "%s"}`, l.UserName, l.Password, l.ApiKey)

	reqTokenURI := SetUrl("/v3/login")

	if l.getCache(l.createKey(reqTokenURI, body)) {
		return l
	}

	_, l.Body, errs = gr.New().
		Post(reqTokenURI).
		Type("form").
		Send(body).
		End()

	for _, err := range errs {
		l.Error = err
	}

	if err := json.Unmarshal([]byte(l.Body), &l.Response); err != nil {
		l.Error = err
	} else {
		isCacheOn := l.OnCache
		l.setCache(l.createKey(reqTokenURI, body), l.Body)
		if isCacheOn {
			l.OnCache = true
		}

	}
	l.SetToken()
	return l
}

func NewLoketApi(url, username, password, key string) (*Loket, error) {
	baseUrl = url
	l := &Loket{
		UserName:     username,
		Password:     password,
		ApiKey:       key,
		Token:        "",
		TokenExpired: true,
	}
	return l, nil
}

func SetUrl(url string) string {
	t := fmt.Sprintf("%s%s", baseUrl, url)
	return t
}

func (l *Loket) SetToken() *Loket {
	resp := struct {
		Status string `json:"status"`
		Data   *struct {
			Token string `json:"token"`
		} `json:"data"`
		Message string `json:"message"`
	}{"", nil, ""}
	byt := []byte(l.Body)

	if err := json.Unmarshal(byt, &resp); err != nil {
		return l
	}

	l.Token = resp.Data.Token
	l.TokenExpired = false
	return l
}

func (l *Loket) SetStruct(v interface{}) *Loket {
	err := json.Unmarshal([]byte(l.Body), &v)
	if err != nil {
		l.Error = err
		return l
	}
	return l
}

func (l *Loket) Post(url, t, body string) *Loket {
	var errs []error
	var err error

	aURL := SetUrl(url)

	if l.getCache(l.createKey(aURL, t, body)) {
		return l
	}

	_, l.Body, errs = gr.New().
		Post(aURL).
		Set("token", l.Token).
		Type(t).
		Send(body).
		End()

	for _, err = range errs {
		l.Error = err
	}

	if err = json.Unmarshal([]byte(l.Body), &l.Response); err != nil {
		l.Error = err
	}

	if len(l.Response.Data.([]interface{})) == 0 {
		l.getAuth().Post(url, t, body)
	} else {
		l.setCache(l.createKey(aURL, t, body), l.Body)

	}

	return l
}

func (l *Loket) Get(url string) *Loket {
	var errs []error
	var err error

	aURL := SetUrl(url)

	if l.getCache(l.createKey(aURL)) {
		return l
	}

	_, l.Body, errs = gr.New().
		Get(aURL).
		Set("token", l.Token).
		End()

	for _, err = range errs {
		l.Error = err
	}

	if err = json.Unmarshal([]byte(l.Body), &l.Response); err != nil {
		l.Error = err
	}

	if len(l.Response.Data.([]interface{})) == 0 {
		l.getAuth().Get(url)
	} else {
		l.setCache(l.createKey(aURL), l.Body)
	}

	return l
}

func (l *Loket) createKey(keys ...string) string {
	for k, v := range keys {
		keys[k] = strings.ToLower(v)
	}
	s := md5.Sum([]byte(strings.Join(keys, "")))
	return fmt.Sprintf("%x", string(s[:]))
}

func (l *Loket) setCache(key, data string) {
	if l.OnCache {
		fmt.Println("Setting cache")
		cache := librarian.Get("redis.master").(*cache.CRedis)

		cache.Set(key, data, CACHE_DURATION)
		l.OnCache = false
	}
}

func (l *Loket) getCache(key string) bool {
	if l.OnCache {
		fmt.Println("Getting cache")
		cache := librarian.Get("redis.slave").(*cache.CRedis)

		c := cache.Get(key)

		if c != "" {
			l.OnCache = false
			l.Body = c
			if err := json.Unmarshal([]byte(l.Body), &l.Response); err != nil {
				l.Error = err
			}
			return true
		}
	}
	return false
}

func (l *Loket) CacheOn() *Loket {
	l.OnCache = true
	return l
}
