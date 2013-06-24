package main

import (
	"github.com/Craig-Macomber/election/msg"
	"github.com/Craig-Macomber/election/server"
	"github.com/Craig-Macomber/election/server/ballot"
	"github.com/Craig-Macomber/election/server/vote"
)

const maxLength = 4096

func main() {
	handlers := server.HandlerMap{
		msg.SignatureRequest: server.BlockHandler(ballot.HandelSignatureRequest, maxLength),
		msg.Vote:             server.BlockHandler(vote.HandelSignatureRequest, maxLength),
	}
	server.Start(handlers)
}
