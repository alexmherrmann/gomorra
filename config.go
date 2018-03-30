package gomorra

import (
	"io/ioutil"
	"encoding/json"
	"errors"
)

// If both privatekeypath and password are specified it indicates that the private key has a password
type HostConfig struct {
	Hostname       string
	PrivateKeyPath string
	Password       string
	Username       string

	// The display name in the dashboard
	Prettyname string
}

type Config struct {
	Revision int
	Hosts    []HostConfig
}

func ReadConfigFile(path string) (Config, error) {
	bytes, err := ioutil.ReadFile(path)

	config := Config{}
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)

	if err != nil {
		return config, err
	}
	if config.Revision != 1 {
		return config, errors.New("Only revision 1 supported")
	}

	return config, nil
}
