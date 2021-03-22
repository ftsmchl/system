package client

func (c *Client) HostFindContracts() error {
	err := c.get("/findContracts", nil)
	return err
}

func (c *Client) HostRegister() error {
	err := c.get("/hostRegister", nil)
	return err
}
