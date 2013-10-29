package config

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg/msgs"
)

// TODO: move this elsewhere
var electionPath string = "demoElection/"

var Path string = electionPath + "config"

func LoadBytes() []byte {
	return keys.LoadBytes(Path)
}

func Unpack(data []byte) *msgs.ElectionConfig {
	var config msgs.ElectionConfig
	err := proto.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &config
}

func Load() *msgs.ElectionConfig {
	return Unpack(LoadBytes())
}

// TODO: move this elsewhere
func LoadServerKey(name string) *msgs.PrivateKey {
	return keys.LoadPrivateKey(electionPath + "serverPrivate/" + name)
}
