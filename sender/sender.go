package sender

import "github.com/thetayloredman/mfc/crypto/jsonsigning"

type Sender struct {
	identity jsonsigning.SigningKey
}

func NewSender(identity jsonsigning.SigningKey) *Sender {
	return &Sender{
		identity: identity,
	}
}
