package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//A client makes requests to the siad HTTP API
type Client struct {
	//Address is the API adddress of the sysd server
	Address string
}

//New creates a client using the provided address
func New(address string) *Client {
	return &Client{
		Address: address,
	}
}

func drainAndClose(rc io.ReadCloser) {
	io.Copy(ioutil.Discard, rc)
	rc.Close()
}

//NewRequest constructs a request to the systemd HTTP API
//The resource path must begin with /.
func (c *Client) NewRequest(method, resource string, body io.Reader) (*http.Request, error) {
	url := "http://" + c.Address + resource
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *Client) get(resource string, obj interface{}) error {
	resp, err := c.getRawResponse(resource)
	if err != nil {
		return err
	}
	fmt.Println(string(resp))
	return nil
}

func (c *Client) getRawResponse(resource string) ([]byte, error) {
	req, err := c.NewRequest("GET", resource, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer drainAndClose(res.Body)
	return ioutil.ReadAll(res.Body)
}
