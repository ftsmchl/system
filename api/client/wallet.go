package client

import (
	"fmt"
)

func (c *Client) WalletAddAccount(publicKey, privateKey string) error {
	fmt.Println("Inside httpClient.WalletAddAccount")
	source := "/addAccount/" + publicKey + "/" + privateKey
	err := c.get(source, nil)
	return err
}
