package config

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg/msgs"
)

// TODO: move this elsewhere
var electionPath string = "demoElection/"

var Path string = electionPath + "config"

func Load() *msgs.ElectionConfig {
	data := keys.LoadBytes(Path)
	var config msgs.ElectionConfig
	err := proto.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &config
}

// TODO: move this elsewhere
func LoadServerKey(name string) *msgs.PrivateKey {
	return keys.LoadPrivateKey(electionPath + "serverPrivate/" + name)
}
