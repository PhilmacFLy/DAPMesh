package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml"
)

type GossiperConf struct {
	BindingPort          int
	Members              []string
	RetransmitMultiplier int
}

type HTTPConf struct {
	Port int
	TLS  bool
	Key  string
	Cert string
}

type AuthConf struct {
	Enabled  bool
	Database string
}

//Config is the struct to save and load the config file
type Config struct {
	Gossiper       GossiperConf
	HTTP           HTTPConf
	Authentication AuthConf
}

//LoadConfig loads the config from the filepath and gives back a Config or an error
func LoadConfig(filepath string) (Config, error) {
	var res Config
	file, err := os.Open(filepath)
	if err != nil {
		return res, errors.New("Error opening file: " + err.Error())
	}
	defer file.Close()
	decoder := toml.NewDecoder(file)
	err = decoder.Decode(&res)
	if err != nil {
		return res, errors.New("Error decoding file: " + err.Error())
	}
	return res, nil

}
