package gomorra

import (
	"golang.org/x/crypto/ssh"
	"fmt"
	"io/ioutil"
	"bytes"
	"log"
)

type Remote struct {
	// This is the Hostname of the remote and is where calls will be made
	Hostname string

	// The ssh client associated with this remote, check that it's not nil before using
	client *ssh.Client

	// This will be nil until getcores is run for the first time
	cores *int

	// This will be nil until the first time it's checked
	totalMemKb *int
}

/*
 * Open up connection using the current users private key file
 */
func (r *Remote) Open(username string, privatekeypath string) error {

	privateBytes, err := ioutil.ReadFile(privatekeypath)
	FatalErr(err)
	signer, err := ssh.ParsePrivateKey(privateBytes)
	FatalErr(err)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		// TODO: Change the below to something more secure in time
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	//config.SetDefaults()

	createdClient, err := ssh.Dial("tcp", r.Hostname, config)
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed connecting to %s", r.Hostname)
	}

	r.client = createdClient

	return nil
}

// Just a little function to help with testing
func (r *Remote) LsDir(path string) string {
	sesh, err := r.client.NewSession()
	FatalErr(err)

	var stdoutBuf bytes.Buffer
	sesh.Stdout = &stdoutBuf
	err = sesh.Run(fmt.Sprintf("ls %s", path))
	FatalErr(err)

	return stdoutBuf.String()
}

func (r *Remote) Close() {
	r.client.Close()
}
