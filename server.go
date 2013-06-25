package main

import (
	"github.com/Craig-Macomber/election/msg"
	"github.com/Craig-Macomber/election/server"
	"github.com/Craig-Macomber/election/server/ballot"
	"github.com/Craig-Macomber/election/server/vote"
	"github.com/Craig-Macomber/election/server/voterlist"
)

const maxLength = 4096

func main() {
	handlers := server.HandlerMap{
		msg.SignatureRequest:    server.BlockHandler(ballot.HandelSignatureRequest, maxLength),
		msg.Vote:                server.BlockHandler(vote.HandelSignatureRequest, maxLength),
		msg.KeySignatureRequest: server.BlockHandler(voterlist.HandelSignatureRequest, maxLength),
	}
	server.Start(handlers)
}
