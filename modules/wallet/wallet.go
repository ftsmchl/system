package wallet

type Wallet struct {
	//ethereum account of our node
	accountAddress string
}

//creation of an empty wallet
func New() *Wallet {
	return &Wallet{}
}

//sets the account we are going to use
func (w *Wallet) SetPrimaryAccount(account string) {
	w.accountAddress = account
}

func (w *Wallet) GetPrimaryAccount() string {
	return w.accountAddress
}
