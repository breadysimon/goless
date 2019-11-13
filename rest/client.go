package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/breadysimon/goless/logging"
	"github.com/breadysimon/goless/util"
	"github.com/gertd/go-pluralize"
	"gopkg.in/resty.v1"
)

var log *logging.Logger = logging.GetLogger()

type Client struct {
	url         string
	api         *resty.Client
	token       string
	err         error
	result      interface{}
	errorResult interface{}
}

func NewClient(baseUrl, rootPath, username, password string) *Client {
	c := &Client{url: baseUrl + rootPath}
	c.api = resty.New()
	return c.Login(username, password)
}
func (c *Client) Error() error {
	return c.err
}

type token struct {
	Code   int    `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

func (c *Client) Login(username, password string) *Client {
	if c.err == nil {
		var u *url.URL
		if u, c.err = url.Parse(c.url); c.err == nil {
			url := fmt.Sprintf("%s://%s:%s/login", u.Scheme, u.Hostname(), u.Port())

			var t token
			var resp *resty.Response
			if resp, c.err = c.api.R().
				SetFormData(map[string]string{
					"username": username,
					"password": password,
				}).
				SetResult(&t).
				Post(url); c.err == nil {
				if resp.StatusCode() == http.StatusOK {
					c.token = t.Token
					log.Info("login success: ", username)
					log.Debug(c.token)
				} else {
					c.err = fmt.Errorf("login failed: %s", username)
				}
			}
		}
	}
	return c
}

type ErrorResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Client) prepareRequest(body, result, errorResult interface{}) *resty.Request {
	r := c.api.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.token)
	if body != nil {
		r = r.SetBody(body)
	}
	if result != nil {
		r = r.SetResult(result)
	}
	if errorResult != nil {
		r = r.SetError(errorResult)
	}
	return r
}

func (c *Client) checkError(err error, e *ErrorResult, resp *resty.Response, expected int) bool {
	switch {
	case err != nil:
		c.err = err
		c.errorResult = nil
		return false
	case resp.StatusCode() != expected:
		c.errorResult = e
		if e.Code == 0 {
			e.Code = resp.StatusCode()
			e.Message = resp.Status()
		}
		c.err = fmt.Errorf("failed. [%s]%s: %v", resp.Request.Method, resp.Request.URL, e)
		return false
	default:
		return true
	}
}

func (c *Client) httpPut(path string, body interface{}, result interface{}) *Client {
	var e ErrorResult
	resp, err := c.prepareRequest(body, result, &e).Put(c.url + path)

	if c.checkError(err, &e, resp, http.StatusOK) {
		c.result = result
	}
	return c
}
func (c *Client) httpPost(path string, body interface{}, result interface{}) *Client {
	var e ErrorResult
	resp, err := c.prepareRequest(body, result, &e).Post(c.url + path)

	if c.checkError(err, &e, resp, http.StatusOK) {
		c.result = result
	}
	return c
}

func (c *Client) httpGet(path string, result interface{}) *Client {
	var e ErrorResult
	resp, err := c.prepareRequest(nil, result, &e).Get(c.url + path)

	if c.checkError(err, &e, resp, http.StatusOK) {
		c.result = result
	}

	return c
}
func (c *Client) httpDelete(path string) *Client {
	var e ErrorResult
	resp, err := c.prepareRequest(nil, nil, &e).Delete(c.url + path)

	c.checkError(err, &e, resp, http.StatusOK)
	return c
}
func guessResourceName(obj interface{}) string {
	tp := reflect.TypeOf(obj).String()
	reg := regexp.MustCompile(`.*\.`)
	name := reg.ReplaceAllString(tp, "")
	pluralize := pluralize.NewClient()
	return strings.ToLower(pluralize.Plural(name))
}
func (c *Client) Create(body interface{}) *Client {
	if c.err == nil {
		reflect.TypeOf(body)
		var result interface{}
		return c.httpPost(guessResourceName(body), body, result)
	}
	return c
}
func (c *Client) Update(body interface{}) *Client {
	if c.err == nil {
		reflect.TypeOf(body)
		var result interface{}
		url := fmt.Sprintf("%s/%d", guessResourceName(body), GetId(body).Int())
		return c.httpPut(url, body, result)
	}
	return c
}

func GetId(m interface{}) reflect.Value {
	v := reflect.Indirect(reflect.ValueOf(m))
	return v.FieldByName("ID")
}

func (c *Client) Find(result interface{}) *Client {
	if c.err == nil {
		resultType := reflect.TypeOf(result).String()
		if util.StartWith(resultType, "*[]") {
			return c.httpGet(guessResourceName(result), result)
		} else {
			url := fmt.Sprintf("%s/%d", guessResourceName(result), GetId(result).Int())
			return c.httpGet(url, result)
		}
	}
	return c
}
func (c *Client) Delete(obj interface{}) *Client {
	if c.err == nil {
		url := fmt.Sprintf("%s/%d", guessResourceName(obj), GetId(obj).Int())
		return c.httpDelete(url)
	}
	return c
}
func (c *Client) peekResult() *Client {
	if c.err == nil {
		fmt.Println("=========")
		fmt.Println(c.result)
	}
	return c
}
