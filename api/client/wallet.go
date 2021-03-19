package client

import ()

func (c *Client) WalletSetAccount(acc string) error {
	err := c.get("/setAccount/address/"+acc, nil)
	return err
}
