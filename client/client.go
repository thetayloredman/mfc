package client

import (
	"github.com/thetayloredman/mfc/crypto/jsonsigning"
	"github.com/thetayloredman/mfc/sender"
)

type Client struct {
	Sender *sender.Sender
}

func NewClient(identity jsonsigning.SigningKey) *Client {
	return &Client{
		Sender: sender.NewSender(identity),
	}
}
