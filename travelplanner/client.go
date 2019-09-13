package travelplanner

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	ErrStatusCode = fmt.Errorf("bad status code")
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	target, err := url.Parse(c.baseURL + path + "&key=" + c.apiKey)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		if rw, ok := body.(io.ReadWriter); ok {
			buf = rw
		} else if str, ok := body.(string); ok {
			buf = bytes.NewBuffer([]byte(str))
		} else if b, ok := body.([]byte); ok {
			buf = bytes.NewBuffer(b)
		} else {
			buf = new(bytes.Buffer)
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			err := enc.Encode(body)
			if err != nil {
				return nil, err
			}
		}
	}

	req, err := http.NewRequest(method, target.String(), buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) DoRaw(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		return resp, ErrStatusCode
	}

	return resp, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.DoRaw(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := xml.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

type Config struct {
	HttpClient *http.Client
	APIKey     string
	BaseURL    string
}

func NewClient(config *Config) *Client {
	if config == nil {
		config = &Config{}
	}
	c := &Client{
		httpClient: config.HttpClient,
		apiKey:     config.APIKey,
		baseURL:    config.BaseURL,
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	if c.apiKey == "" {
		c.apiKey = os.Getenv("TRAFIKLAB_TRAVELPLANNER13_API_KEY")
	}
	if c.baseURL == "" {
		c.baseURL = os.Getenv("TRAFIKLAB_TRAVELPLANNER31_BASE_URL")
		if c.baseURL == "" {
			c.baseURL = "http://api.sl.se"
		}
	}

	return c
}
