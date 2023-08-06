package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg/docker"
	"github.com/ustclug/podzol/pkg/server"
)

type BadStatusCodeError struct {
	StatusCode int
	Message    string
}

func (e BadStatusCodeError) Error() string {
	return fmt.Sprintf("bad status code %d: %s", e.StatusCode, e.Message)
}

// Client is a client for the podzol server.
type Client struct {
	serverAddr string
	httpClient *http.Client
}

// NewClient creates a new client from config.
func NewClient(v *viper.Viper) *Client {
	return &Client{
		serverAddr: v.GetString("listen-addr"),
		httpClient: &http.Client{
			Timeout: v.GetDuration("timeout"),
		},
	}
}

// MakeURL creates a URL from the given path.
func (c *Client) MakeURL(path string) string {
	return "http://" + c.serverAddr + path
}

// MakeRequest constructs a *http.Request from the given arguments.
func (c *Client) MakeRequest(method, path string, payload any) (*http.Request, error) {
	body := new(bytes.Buffer)
	if payload != nil {
		if err := json.NewEncoder(body).Encode(payload); err != nil {
			return nil, err
		}
	}
	return http.NewRequest(method, c.MakeURL(path), body)
}

// doRequest performs a request and decodes the response into output.
func (c *Client) doRequest(method, path string, input, output any) error {
	req, err := c.MakeRequest(method, path, input)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Attempt to decode error message
		var errResp server.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return BadStatusCodeError{StatusCode: resp.StatusCode, Message: errResp.Error}
		}
		// Decode failed, message unavailable
		return BadStatusCodeError{StatusCode: resp.StatusCode}
	}

	// do not attempt to decode if output is not required
	if output == nil {
		return nil
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		return fmt.Errorf("bad content type: %s", contentType)
	}

	return json.NewDecoder(resp.Body).Decode(output)
}

func (c *Client) Create(opts docker.ContainerOptions) (data docker.ContainerInfo, err error) {
	err = c.doRequest(http.MethodPost, "/create", opts, &data)
	return
}

func (c *Client) Remove(opts docker.ContainerOptions) (err error) {
	err = c.doRequest(http.MethodPost, "/remove", opts, nil)
	return
}

func (c *Client) List(opts docker.ContainerOptions) (data []docker.ContainerInfo, err error) {
	err = c.doRequest(http.MethodPost, "/list", opts, &data)
	return
}
