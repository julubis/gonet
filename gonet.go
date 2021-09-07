package gonet

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Request struct {
	httpreq *http.Request
	Header  *http.Header
	Client  *http.Client
	Cookies []*http.Cookie
}

type Response struct {
	R       *http.Response
	content []byte
	text    string
	req     *Request
}

type Header map[string]string

func Session() *Request {
	req := new(Request)
	req.httpreq = &http.Request{
		Method: "GET",
		Header: make(http.Header),
		Proto: "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	req.Header = &req.httpreq.Header
	req.Client = &http.Client{}

	jar, _ := cookiejar.New(nil)
	req.Client.Jar = jar
	return req
}

func (req *Request) SetCookie(cookie *http.Cookie) {
	req.Cookies = append(req.Cookies, cookie)
}

func (req *Request) SetUrl(Url string) {
	parseUrl, _ := url.Parse(Url)
	req.httpreq.URL = parseUrl
}

func (req *Request) SetBody(Body string) {
	req.httpreq.Body = ioutil.NopCloser(strings.NewReader(Body))
	req.Header.Set("Content-Type", "application/json")
}

func (req *Request) Gets() (resp *Response) {
	req.httpreq.Method = "GET"
	delete(req.httpreq.Header, "Cookie")
	req.Client.Jar.SetCookies(req.httpreq.URL, req.Cookies)
	res, _ := req.Client.Do(req.httpreq)
	resp = &Response{}
	resp.R = res
	resp.req = req
	resp.Content()
	defer res.Body.Close()
	return resp
}

func (req *Request) Posts() (resp *Response) {
	req.httpreq.Method = "POST"
	delete(req.httpreq.Header, "Cookie")

	req.Client.Jar.SetCookies(req.httpreq.URL, req.Cookies)
	res, _ := req.Client.Do(req.httpreq)

	resp = &Response{}
	resp.R = res
	resp.req = req

	resp.Content()
	defer res.Body.Close()
	return resp
}

func Get(Url string) (resp *Response) {
	req := new(Request)
	req.httpreq = &http.Request{
		Method: "GET",
		Header: make(http.Header),
		Proto: "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	req.Header = &req.httpreq.Header
	req.Client = &http.Client{}

	jar, _ := cookiejar.New(nil)
	req.Client.Jar = jar

	delete(req.httpreq.Header, "Cookie")
	parseUrl, _ := url.Parse(Url)
	req.httpreq.URL = parseUrl

	res, _ := req.Client.Do(req.httpreq)
	resp = &Response{}
	resp.R = res
	resp.req = req
	resp.Content()
	defer res.Body.Close()
	return resp
}

func Post(Url string, Body string) (resp *Response) {
	req := new(Request)
	req.httpreq = &http.Request{
		Method: "POST",
		Header: make(http.Header),
		Proto: "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	req.Header = &req.httpreq.Header
	req.Client = &http.Client{}

	jar, _ := cookiejar.New(nil)
	req.Client.Jar = jar

	req.Header.Set("Content-Type", "application/json")
	delete(req.httpreq.Header, "Cookie")

	req.httpreq.Body = ioutil.NopCloser(strings.NewReader(Body))
	URL, _ := url.Parse(Url)
	req.httpreq.URL = URL

	res, _ := req.Client.Do(req.httpreq)
	// clear post  request information
	req.httpreq.Body = nil
	req.httpreq.GetBody = nil
	req.httpreq.ContentLength = 0

	resp = &Response{}
	resp.R = res
	resp.req = req

	resp.Content()
	defer res.Body.Close()
	return resp
}

// *** kumpulan method request ***
func (req *Request) Get(Url string) (resp *Response) {
	req.httpreq.Method = "GET"
	delete(req.httpreq.Header, "Cookie")
	parseUrl, _ := url.Parse(Url)
	req.httpreq.URL = parseUrl

	req.Client.Jar.SetCookies(req.httpreq.URL, req.Cookies)
	res, _ := req.Client.Do(req.httpreq)
	resp = &Response{}
	resp.R = res
	resp.req = req
	resp.Content()
	defer res.Body.Close()
	return resp
}
func (req *Request) Post(Url string, Body string) (resp *Response) {
	req.httpreq.Method = "POST"
	req.Header.Set("Content-Type", "application/json")
	delete(req.httpreq.Header, "Cookie")

	req.httpreq.Body = ioutil.NopCloser(strings.NewReader(Body))
	URL, _ := url.Parse(Url)
	req.httpreq.URL = URL

	req.Client.Jar.SetCookies(req.httpreq.URL, req.Cookies)
	res, _ := req.Client.Do(req.httpreq)

	// clear post  request information
	req.httpreq.Body = nil
	req.httpreq.GetBody = nil
	req.httpreq.ContentLength = 0

	resp = &Response{}
	resp.R = res
	resp.req = req

	resp.Content()
	defer res.Body.Close()
	return resp
}

func (resp *Response) Content() []byte {
	var err error
	if len(resp.content) > 0 {
		return resp.content
	}
	var Body = resp.R.Body
	if resp.R.Header.Get("Content-Encoding") == "gzip" && resp.req.Header.Get("Accept-Encoding") != "" {
		reader, err := gzip.NewReader(Body)
		if err != nil {
			return nil
		}
		Body = reader
	}
	resp.content, err = ioutil.ReadAll(Body)
	if err != nil {
		return nil
	}
	return resp.content
}

// *** kumpulan type response ***
func (resp *Response) Text() string {
	if resp.content == nil {
		resp.Content()
	}
	resp.text = string(resp.content)
	return resp.text
}
func (resp *Response) Json(v interface{}) error {
	if resp.content == nil {
		resp.Content()
	}
	return json.Unmarshal(resp.content, v)
}
func (resp *Response) Cookies() (cookies []*http.Cookie) {
	httpreq := resp.req.httpreq
	client := resp.req.Client
	cookies = client.Jar.Cookies(httpreq.URL)
	return cookies
}