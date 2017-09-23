package ssh

import (
	//"bufio"
	//"fmt"
	"os"
	"os/exec"
	//"time"
	"io"
	"github.com/paulczar/pty"
)

// SSH settings
type Options struct {
	UserName  		string
	Password		  string
	Host					string
	Port      		string
	SSHConfigFile	string
}

// SSH into server
func SSH(options *Options) error {
	var (
		sshArgs []string
	)

	if options.SSHConfigFile != "" {
		sshArgs = append(sshArgs,
			[]string{"-F", options.SSHConfigFile}...)
	}

	sshArgs = append(sshArgs,
		[]string{options.Host}...)
	cmd := exec.Command("ssh", sshArgs...)
	
	f, restore, err := pty.StartRaw(cmd)
	if err != nil {
			return err
	}

	defer restore()

	go func() {
			io.Copy(f, os.Stdin)
	}()
	io.Copy(os.Stdout, f)	

  return nil
}