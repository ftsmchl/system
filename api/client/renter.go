package client

import ()

func (c *Client) RenterCreateContracts() error {
	err := c.get("/createAuction", nil)
	return err
}

func (c *Client) RenterUploadFile(source string) error {
	err := c.get("/uploadFile/"+source, nil)
	return err
}
