package spark

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/romana/rlog"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RESTClient interface {
	Get(url string, params map[string]string, out interface{}) error
	Delete(url string, params map[string]string, out interface{}) error
	Post(url string, body interface{}, out interface{}) error
}

type restClient struct {
	http  HTTPClient
	url   string
	token string
}

func NewRESTClient(url, token string) RESTClient {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return &restClient{
		http:  http.DefaultClient,
		url:   url,
		token: token,
	}
}

// err forms an error message with the spark token masked
func (c *restClient) err(format string, args ...interface{}) error {
	str := fmt.Sprintf(format, args...)
	str = strings.Replace(str, c.token, "REDACTED", -1)
	return errors.New(str)
}

// dbg writes a debug log with the spark token masked
func (c *restClient) dbg(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	str = strings.Replace(str, c.token, "REDACTED", -1)
	rlog.Debug(str)
}

// do performs a request forms an error message with the spark token masked
func (c *restClient) do(req *http.Request, out interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	c.dbg("Making request %+v", req)
	resp, err := c.http.Do(req)
	c.dbg("Got response %+v", resp)

	if err != nil {
		return c.err("Request %+v failed: %v", req, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return c.err("Request %+v had an error response: %+v", req, resp)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			rlog.Info("Error closing response body: %v")
		}
	}()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.err("Request %+v could not read body from %+v: %v", req, resp, err)
	}
	c.dbg("Got response body %s", bytes)

	if err = json.Unmarshal(bytes, &out); err != nil {
		return c.err("Could not unmarshal %s into %+v: %v", bytes, out, err)
	}
	return nil
}

func (c *restClient) Get(path string, params map[string]string, out interface{}) error {
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}

	url, err := url.Parse(c.url + path)
	if err != nil {
		return err
	}
	url.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *restClient) Delete(path string, params map[string]string, out interface{}) error {
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}

	url, err := url.Parse(c.url + path)
	if err != nil {
		return err
	}
	url.RawQuery = query.Encode()

	req, err := http.NewRequest("DELETE", url.String(), nil)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *restClient) Post(path string, body interface{}, out interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.url+path, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	return c.do(req, out)
}
