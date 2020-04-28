package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	timeout  = 3 * time.Second
	basePath = "http://169.254.169.254"
)

type Client struct {
	client  *http.Client
	baseUrl *url.URL
}

func NewClient() *Client {

	u, err := url.Parse(basePath)
	if err != nil {
		panic(err)
	}

	c := &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		baseUrl: u,
	}

	return c
}

// Metadata returns the entire contents of the instances metadata
func (c *Client) Metadata() (*MetaData, error) {
	metadata := &MetaData{}

	err := c.doRequest("/v1.json", metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (c *Client) doRequest(uri string, meta *MetaData) error {

	resp, err := c.client.Get(fmt.Sprintf("%s%s", c.baseUrl, uri))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, meta); err != nil {
		return err
	}

	return nil
}

func (c *Client) SetBaseURL(baseUrl string) error {
	updatedURL, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}

	c.baseUrl = updatedURL
	return nil
}
