package wallet

type Wallet struct {
	//ethereum account of our node
	pubKey  string
	privKey string
}

//creation of an empty wallet
func New() *Wallet {
	return &Wallet{}
}

//sets the account we are going to use
func (w *Wallet) SetPrimaryAccount(publicKey, privateKey string) {
	w.pubKey = publicKey
	w.privKey = privateKey
}

func (w *Wallet) GetPrimaryAccount() string {
	return w.pubKey
}

func (w *Wallet) GetPublicKey() string {
	return w.pubKey
}

func (w *Wallet) GetPrivateKey() string {
	return w.privKey
}
