package client

func (c *Client) HostFindContracts() error {
	err := c.get("/findContracts", nil)
	return err
}
