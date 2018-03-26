package gomorra

import (
	"bytes"
	"errors"
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
