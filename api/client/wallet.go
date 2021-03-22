package client

import (
	"fmt"
)

func (c *Client) WalletAddAccount(acc string) error {
	fmt.Println("Inside httpClient.WalletAddAccount")
	err := c.get("/addAccount/"+acc, nil)
	return err
}
