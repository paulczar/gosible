package ansible

import (
//  	"fmt"
//    "strings"
)

// PlaybookOptions for running ansible
type PlaybookOptions struct {
}

// PlaybookRun ansible playbook
func PlaybookRun(options *Options, ansibleArgs []string) error {
  var (
		err         error
    gosibleArgs []string
	)

  err = configureEnvironment(options)
  if err != nil {
    return err
  }
  err = sshConfigFile(options)
  if err != nil {
    return err
  }
  err = configureKnownHostsFile(options)
  if err != nil {
    return err
  }
  configureSSHForwardAgent(options)

  if options.Inventory != "" {
    gosibleArgs = append(gosibleArgs,
      []string{"--inventory", options.Inventory}...)
  }

  cmdName := "ansible-playbook"
	cmdArgs := append(gosibleArgs, ansibleArgs...)
  //fmt.Println("running: ansible_playbook", strings.Join(cmdArgs, " "))

  return runCmd(cmdName, cmdArgs)
}
