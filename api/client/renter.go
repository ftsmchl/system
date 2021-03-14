package client

import ()

func (c *Client) RenterCreateContracts() error {
	err := c.get("/createAuction", nil)
	return err
}
