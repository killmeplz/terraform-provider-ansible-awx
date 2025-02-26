package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	Host   string
	Token  string
	Client *http.Client
}

func NewClient(host, token string) (*Client, error) {
	if host == "" || token == "" {
		return nil, errors.New("host and token must be provided")
	}
	return &Client{
		Host:   host,
		Token:  token,
		Client: &http.Client{},
	}, nil
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.Host, path)
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	return req, nil
}

func (c *Client) do(req *http.Request) (map[string]interface{}, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for non-200 status codes
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result map[string]interface{}
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (c *Client) Post(path string, body interface{}) (map[string]interface{}, error) {
	req, err := c.newRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) Get(path string) (map[string]interface{}, error) {
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) Patch(path string, body interface{}) (map[string]interface{}, error) {
	req, err := c.newRequest("PATCH", path, body)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) Put(path string, body interface{}) (map[string]interface{}, error) {
	req, err := c.newRequest("PUT", path, body)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) Delete(path string) error {
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req)
	return err
}

func (c *Client) IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "HTTP 404")
}
