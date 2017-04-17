package base

import (
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestSshConnection(t *testing.T) {
	s := Server{}
	if err := s.configFromFile("./config.yml"); err != nil {
		t.Error(err)
	}
	client, err := ssh.Dial("tcp", s.sshServer, s.sshConfig)
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		t.Error(err)
	}
	res, err := sess.CombinedOutput("hostname")
	t.Logf("%s", res)
}
