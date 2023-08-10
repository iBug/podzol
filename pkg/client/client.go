package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

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
	verbose    bool
}

// NewClient creates a new client from config.
func NewClient(v *viper.Viper) *Client {
	return &Client{
		serverAddr: v.GetString("listen-addr"),
		httpClient: &http.Client{
			Timeout: v.GetDuration("timeout"),
		},
		verbose: v.GetBool("verbose"),
	}
}

// MakeURL creates a URL from the given path.
func (c *Client) makeURL(path string) string {
	return "http://" + c.serverAddr + path
}

// MakeRequest constructs a *http.Request from the given arguments.
func (c *Client) makeRequest(method, path string, payload any) (*http.Request, error) {
	url := c.makeURL(path)
	body := new(bytes.Buffer)
	if payload != nil {
		if err := json.NewEncoder(body).Encode(payload); err != nil {
			return nil, err
		}
	}
	if c.verbose {
		// Produce a copy of the request to stderr
		fmt.Fprintf(os.Stderr, "> %s %s\n", method, url)
		if payload != nil {
			e := json.NewEncoder(os.Stderr)
			e.SetIndent("", "  ")
			if err := e.Encode(payload); err != nil {
				fmt.Fprintf(os.Stderr, "error encoding payload: %s\n", err)
			}
		}
		fmt.Fprintln(os.Stderr)
	}
	return http.NewRequest(method, url, body)
}

// doRequest performs a request and decodes the response into output.
func (c *Client) doRequest(method, path string, input, output any) error {
	req, err := c.makeRequest(method, path, input)
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

	reader := io.Reader(resp.Body)
	if c.verbose {
		pipeR, pipeW := io.Pipe()
		reader = io.TeeReader(resp.Body, pipeW)

		wg := sync.WaitGroup{}
		wg.Add(1)
		defer wg.Wait()
		defer pipeW.Close()
		go func() {
			defer wg.Done()
			var buf any
			if err := json.NewDecoder(pipeR).Decode(&buf); err != nil {
				fmt.Fprintf(os.Stderr, "error decoding response: %v\n", err)
				return
			}
			fmt.Fprintln(os.Stderr) // Separate response from request
			fmt.Fprintf(os.Stderr, "< HTTP %s\n", resp.Status)
			e := json.NewEncoder(os.Stderr)
			e.SetIndent("", "  ")
			if err := e.Encode(buf); err != nil {
				fmt.Fprintf(os.Stderr, "error printing response: %v\n", err)
			}
		}()
	}

	return json.NewDecoder(reader).Decode(output)
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

func (c *Client) Purge() (data []docker.ContainerInfo, err error) {
	err = c.doRequest(http.MethodPost, "/purge", nil, &data)
	return
}
