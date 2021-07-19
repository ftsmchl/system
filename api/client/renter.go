package client

func (c *Client) RenterCreateContracts() error {
	err := c.get("/createAuction", nil)
	return err
}

func (c *Client) RenterUploadFile(source string) error {
	err := c.get("/uploadFile/"+source, nil)
	return err
}

func (c *Client) RenterChallengeHost(hostPublicKey string) error {
	err := c.get("/challengeHost/"+hostPublicKey, nil)
	return err
}
