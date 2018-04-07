package gomorra

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

const envPath string = "/usr/bin/env"

/*
 * Note that this requires an absolute path at the moment
 * It also reads the ENTIRE file into memory so it's not suitable for large files
 */
func (r *Remote) readFileFromSystem(path string) (*bytes.Buffer, error) {
	readBytes := new(bytes.Buffer)
	if r.client == nil {
		return readBytes, errors.New("Client is not open!")
	}

	session, err := r.client.NewSession()
	if err != nil {
		return readBytes, err
	}
	defer session.Close()

	session.Stdout = readBytes
	err = session.Run(envPath + " cat " + path)
	return readBytes, nil
}

// Check that a stat result is an int
func CheckInt(result StatResult) (int, bool) {
	switch v := result.GenericResult.(type) {
	case int:
		return v, true
	default:
		return 0, false
	}
}

func CheckFloat(result StatResult) (float64, bool) {
	switch v := result.GenericResult.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0.0, false
	}
}

func GetSignerFromPrivateKey(privateKeyPath string) (ssh.Signer, error) {
	return nil, NotImplementedErr
}

func GetSignerFromPrivateKeyWithPassword(privateKeyPath string, password string) (ssh.Signer, error) {
	return nil, NotImplementedErr
}

func GetRemoteFromHostConfig(config HostConfig, logger *log.Logger) (*Remote, error) {

	toReturn := new(Remote)
	toReturn.Logger = logger
	toReturn.Hostname = config.Hostname

	if len(config.Username) <= 0 {
		return nil, errors.New("No username in config.json")
	}
	toReturn.username = config.Username

	if len(config.Password) > 0 && len(config.PrivateKeyPath) > 0 {
		return nil, NotImplementedErr
	}

	if len(config.Password) > 0 {
		toReturn.methods = []ssh.AuthMethod{ssh.Password(config.Password)}
	} else if len(config.PrivateKeyPath) > 0 {
		privateBytes, err := ioutil.ReadFile(config.PrivateKeyPath)
		if err != nil {
			return nil, err
		}

		signer, err := ssh.ParsePrivateKey(privateBytes)
		if err != nil {
			return nil, err
		}

		toReturn.methods = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		return nil, errors.New("No auth methods in config.json!")
	}

	return toReturn, nil
}